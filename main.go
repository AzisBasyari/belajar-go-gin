package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type album struct {
	Id     uint8   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{Id: 1, Title: "NAINAINAI", Artist: "Atarashii Gakkou!", Price: 50.00},
	{Id: 2, Title: "Janaindayo", Artist: "Atarashii Gakkou!", Price: 42.50},
	{Id: 3, Title: "Que Sera Sera", Artist: "Atarashii Gakkou!", Price: 48.00},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/album/:id", getAlbumById)
	router.DELETE("/album/:id", deleteAlbumById)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Cannot Process!"})
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")
	int_id, err := strconv.ParseInt(id, 10, 8)

	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Cannot Process!"})
		return
	}

	for _, album := range albums {
		if album.Id == uint8(int_id) {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album Not Found!"})
}

func deleteAlbumById(c *gin.Context) {
	var deleted_idx int

	id := c.Param("id")
	int_id, err := strconv.ParseInt(id, 10, 8)

	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Cannot Process!"})
		return
	}

	for idx, album := range albums {
		if album.Id == uint8(int_id) {
			deleted_idx = idx
		}
	}

	albums = append(albums[:deleted_idx], albums[deleted_idx+1:]...)

	c.IndentedJSON(http.StatusOK, albums)
}
