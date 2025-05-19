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

	router.LoadHTMLGlob("internal/web/templates/*.html")
	router.Static("/static", "internal/web/static")

	api_group := router.Group("/api")
	api.RegisterRoutes(api_group)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
