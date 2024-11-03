package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func initDB() error {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASSWD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	return db.Ping()
}

func fetchAlbums() ([]Album, error) {
	albums := []Album{}
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return albums, err
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return albums, err
		}
		albums = append(albums, alb)
	}

	return albums, nil
}

func fetchAlbumByID(id int64) (album, error) {
	var alb album
	err := db.QueryRow("SELECT * FROM album WHERE id = ?", id).Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	return alb, err
}

func createAlbum(a album) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Printf("statement error: %v", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(a.Title, a.Artist, a.Price)
	if err != nil {
		fmt.Printf("exec error: %v", err)
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("commit error: %v", err)
		return 0, err
	}
	return res.LastInsertId()
}

func fetchAlbumsByArtist(name string) ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist: %q: %v", name, err)
	}
	return albums, nil
}
