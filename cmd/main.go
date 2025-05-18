package main

import (
	"log"
	"weather-forecast/internal/api"
	"weather-forecast/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	router := gin.Default()
	api.RegisterRoutes(router)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
