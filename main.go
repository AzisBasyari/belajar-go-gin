package main

import (
	"database/sql"
	"net/http"

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
	router.PUT("/albums", updateAlbumById)
	router.GET("/album/:id", getAlbumById)
	router.DELETE("/album/:id", deleteAlbumById)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	allAlbums := []album{}

	db := GetConnection()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM main.albums")
	defer rows.Close()

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err})
		return
	}

	for rows.Next() {
		var album album

		err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price)

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err})
			return
		}

		allAlbums = append(allAlbums, album)
	}

	c.JSON(http.StatusOK, allAlbums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err})
		return
	}

	db := GetConnection()
	defer db.Close()

	err := db.QueryRow("INSERT INTO main.albums (title, artist, price) VALUES ($1, $2, $3) RETURNING id", newAlbum.Title, newAlbum.Artist, newAlbum.Price).Scan(&newAlbum.Id)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err})
		return
	}

	err = db.QueryRow("SELECT * FROM main.albums WHERE id = $1", newAlbum.Id).Scan(&newAlbum.Id, &newAlbum.Title, &newAlbum.Artist, &newAlbum.Price)

	switch {
	case err == sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{"message": "Album Not Found!"})
		return
	case err != nil:
		c.JSON(http.StatusForbidden, gin.H{"message": err})
		return
	default:
		c.JSON(http.StatusCreated, newAlbum)
		return
	}
}

func getAlbumById(c *gin.Context) {
	var album album
	id := c.Param("id")

	db := GetConnection()
	defer db.Close()

	err := db.QueryRow("SELECT * FROM main.albums WHERE id = $1", id).Scan(&album.Id, &album.Title, &album.Artist, &album.Price)

	switch {
	case err == sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{"message": "Album Not Found!"})
		return
	case err != nil:
		c.JSON(http.StatusForbidden, gin.H{"message": err})
		return
	default:
		c.JSON(http.StatusOK, album)
		return
	}
}

func deleteAlbumById(c *gin.Context) {
	var album album
	id := c.Param("id")

	db := GetConnection()
	defer db.Close()

	err := db.QueryRow("SELECT * FROM main.albums WHERE id = $1", id).Scan(&album.Id, &album.Title, &album.Artist, &album.Price)

	switch {
	case err == sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{"message": "Album Not Found!"})
		return
	case err != nil:
		c.JSON(http.StatusForbidden, gin.H{"message": err})
		return
	default:
		result, err := db.Exec("DELETE FROM main.albums WHERE id = $1", id)

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err})
			return
		}

		rows, err := result.RowsAffected()

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err})
			return
		}

		if rows != 1 {
			c.JSON(http.StatusForbidden, gin.H{"message": err})
			return
		}

		c.JSON(http.StatusOK, album)
		return
	}
}

func updateAlbumById(c *gin.Context) {
	var updatedAlbum album

	if err := c.BindJSON(&updatedAlbum); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err})
		return
	}

	if updatedAlbum.Id == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Required Id of Album"})
		return
	}

	db := GetConnection()
	defer db.Close()

	err := db.QueryRow("SELECT id FROM main.albums WHERE id = $1", updatedAlbum.Id).Scan(&updatedAlbum.Id)

	switch {
	case err == sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{"message": "Album Not Found!"})
		return
	case err != nil:
		c.JSON(http.StatusForbidden, gin.H{"message": err})
		return
	default:
		result, err := db.Exec("UPDATE main.albums SET title = $1, artist = $2, price = $3 WHERE id = $4", updatedAlbum.Title, updatedAlbum.Artist, updatedAlbum.Price, updatedAlbum.Id)

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err})
			return
		}

		rows, err := result.RowsAffected()

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err})
			return
		}

		if rows != 1 {
			c.JSON(http.StatusForbidden, gin.H{"message": "Expected to affect 1 row"})
			return
		}

		c.JSON(http.StatusOK, updatedAlbum)
	}
}
