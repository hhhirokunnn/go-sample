package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:price`
}

func getAlbums(c *gin.Context) {
	var (
		albums []Album
		err    error
	)

	rawArtist := c.Query("artist")
	if len(rawArtist) > 0 {
		albums, err = fetchAlbumsByArtist(rawArtist)
	} else {
		albums, err = fetchAlbums()
	}

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Printf("newAlbum: %v\n", newAlbum)
	// TODO: validation

	generatedId, err := createAlbum(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	alb, err := fetchAlbumByID(generatedId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, alb)
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
