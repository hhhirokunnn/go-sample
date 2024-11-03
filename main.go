package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}
