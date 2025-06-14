package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hi from gwid",
		})
	})

	log.Printf("Server running on %s\n", ":5000")

	if err := router.Run(":5000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
