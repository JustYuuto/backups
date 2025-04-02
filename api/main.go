package main

import (
	"backups/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	isProd := gin.Mode() == "release"

	if isProd {
		r.Static("/", "./dist")

		r.Group("/api")
		{
			r.GET("/buckets", routes.GetBuckets)
			r.POST("/backup", routes.StartBackup)
		}
	} else {
		r.GET("/buckets", routes.GetBuckets)
		r.POST("/backup", routes.StartBackup)
	}

	r.Run()
}
