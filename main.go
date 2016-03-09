package main

import (
	"net/http"
	"os"

	"github.com/RuniVN/random-girl/models"
	"github.com/jinzhu/gorm"

	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := gorm.Open("postgres", "host="+os.Getenv("POSTGRES_HOST")+" user="+os.Getenv("POSTGRES_USER")+" password="+os.Getenv("POSTGRES_PASSWORD")+" dbname=random-girl sslmode=disable")
	if err != nil {
		panic(err)
	}

	if !db.HasTable(&models.Image{}) {
		db.CreateTable(&models.Image{})
	}

	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/images", func(c *gin.Context) {
			type Request struct {
				Url string `json:"url"`
			}

			var body Request

			if c.BindJSON(&body) == nil {
				if body.Url == "" {
					c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
				}

				var count int

				db.Where("url = ?", body.Url).Count(&count)

				if count > 0 {
					c.JSON(http.StatusFound, gin.H{"status": "this is image has been uploaded"})
				}

				var image models.Image
				image.Url = body.Url
				db.Create(&image)

				c.JSON(http.StatusAccepted, gin.H{"status": "created"})
			}
		})

		v1.GET("/images", func(c *gin.Context) {
			var image models.Image
			err := db.Raw("SELECT * FROM images ORDER BY random() LIMIT 1").Scan(&image)
			if err != nil {
				fmt.Println(err)
			}
			c.JSON(http.StatusOK, image)
		})
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(":8082")
}
