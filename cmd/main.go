package main

import (
	"log"
	"weather-forecast/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	if err := router.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
