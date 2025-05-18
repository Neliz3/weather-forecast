package api

import (
	"net/http"

	"weather-forecast/internal/config"
	"weather-forecast/internal/service"

	"github.com/gin-gonic/gin"
)

func handleGetWeather(c *gin.Context) {
	cfg := config.Load()

	var req struct {
		City string `form:"city" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryParams := map[string]string{
		"q":   req.City,
		"key": cfg.Weather.WEATHER_API_KEY,
	}

	weather, err := service.FetchWeatherNow(cfg.Weather.WEATHER_API_URL, cfg.Weather.WEATHER_API_KEY, queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, weather)
}
