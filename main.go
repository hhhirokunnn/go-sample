package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:price`
}

// var albs = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "hoge", Artist: "fuga", Price: 17.99},
// 	{ID: "3", Title: "piyo", Artist: "foo", Price: 23.99},
// }

var (
	db  *sql.DB
	ctx context.Context
)

func getAlbums(c *gin.Context) {
	albums := []album{}
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		albums = append(albums, alb)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	rawId := c.Param("id")
	id, err := strconv.ParseInt(rawId, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	alb, err := fetchAlbumByID(id)
	switch {
	case err == sql.ErrNoRows:
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	case err != nil:
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	default:
		c.IndentedJSON(http.StatusOK, alb)
	}
}

func fetchAlbumByID(id int64) (album, error) {
	var alb album
	err := db.QueryRow("SELECT * FROM album WHERE id = ?", id).Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	return alb, err
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Printf("newAlbum: %v\n", newAlbum)

	tx, err := db.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Printf("statement error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		fmt.Printf("exec error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("update drivers: unable to rollback: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	generatedId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("generatedID err: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	alb, err := fetchAlbumByID(generatedId)

	c.IndentedJSON(http.StatusCreated, alb)
}

func albumsByArtist(name string) ([]Album, error) {
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

func getAlbumByArtist(c *gin.Context) {
	rawArtist := c.Query("artist")
	album, err := albumsByArtist(rawArtist)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func main() {
	cfg := mysql.Config{
		User: os.Getenv("DBUSER"),
		// Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
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
	fmt.Println("Connected!")
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("albums found: %v\n", albums)

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.GET("/albums_s", getAlbumByArtist)

	router.Run("localhost:8080")
}
