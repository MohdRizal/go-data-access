package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	//setup koneksi
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "golang_recordings",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Terhubung! Urrraaaa!!!")

	albums, err := albumsByArtist("Jono")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	alb, err := albumByID(3)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album Found: %v\n", alb)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID of added album: %v\n", albID)

	delAlb, err := delAlbum(1)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album deleted: %v\n", delAlb)

	albUpdate, err := updateAlbum(Album{
		Title:  "The Modern Sound of Jack Lesmana",
		Artist: "Jack Lesmana",
		Price:  169.99,
	}, 2)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Update album: %v\n", albUpdate)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func albumByID(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumByID %d: no such album", id)
		}
		return alb, fmt.Errorf("albumByID %d: no such album", id)
	}

	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}

func delAlbum(id int64) (int64, error) {
	result, err := db.Exec("DELETE FROM album WHERE id = ?", id)

	if err != nil {
		return 0, fmt.Errorf("delAlbum: %v", err)
	}

	deleteID, err := result.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return deleteID, nil
}

func updateAlbum(alb Album, id int64) (int64, error) {
	result, err := db.Exec("UPDATE album SET title = ?, artist = ?, price = ? WHERE id = ?", alb.Title, alb.Artist, alb.Price, id)

	if err != nil {
		return 0, fmt.Errorf("updateAlbum: %v", err)
	}

	affectedRow, err := result.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("updateAlbum: %v", err)
	}

	return affectedRow, nil
}
