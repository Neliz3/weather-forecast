package main

import (
	"database/sql"
	"fmt"
	"log"
	"weather-forecast/internal/api"
	"weather-forecast/internal/config"
	"weather-forecast/internal/db"
	"weather-forecast/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func startScheduler(cfg *config.Config, dbConn *sql.DB) {
	c := cron.New()

	// Hourly job at minute 0 of every hour
	if _, err := c.AddFunc("0 * * * *", func() {
		sendEmailsByFrequency(dbConn, "hourly", cfg)
	}); err != nil {
		log.Printf("Failed to schedule cron job: %v", err)
	}

	// Daily job at midnight
	if _, err := c.AddFunc("0 0 * * *", func() {
		sendEmailsByFrequency(dbConn, "daily", cfg)
	}); err != nil {
		log.Printf("Failed to schedule cron job: %v", err)
	}

	c.Start()
}

func sendEmailsByFrequency(dbConn *sql.DB, frequency string, cfg *config.Config) {
	subs, err := db.GetConfirmedSubscriptionsByFrequency(dbConn, frequency)
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		return
	}

	for _, sub := range subs {
		queryParams := map[string]string{
			"key": cfg.Weather.API_KEY,
			"q":   sub.City,
		}

		weather, err := service.FetchWeatherNow(cfg.Weather.API_URL, queryParams)
		if err != nil {
			log.Printf("Failed to fetch weather for %s: %v", sub.City, err)
			continue
		}

		err = service.SendWeatherUpdateEmail(cfg.Email.EmailFrom, sub.Email, sub.City, frequency, *weather, cfg)
		if err != nil {
			log.Printf("Failed to send weather email to %s: %v", sub.Email, err)
		}
	}
}

func main() {
	cfg := config.Load()

	dbConn, _ := db.Connect()
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Failed to close DB connection: %v", err)
		}
	}()

	go startScheduler(cfg, dbConn)

	router := gin.Default()

	router.LoadHTMLGlob("internal/web/templates/*.html")
	router.Static("/static", "internal/web/static")

	api_group := router.Group("/api")
	api.RegisterRoutes(api_group)

	if err := router.Run(cfg.Port); err != nil {
		fmt.Printf("Failed to run server: %v", err)
		return
	}
}
