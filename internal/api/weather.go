package api

import (
	"net/http"
	"strings"

	"weather-forecast/internal/config"
	"weather-forecast/internal/service"

	"github.com/gin-gonic/gin"
)

func HandleGetWeather(c *gin.Context) {
	cfg := config.Load()

	var req struct {
		City string `form:"q" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	queryParams := map[string]string{
		"q":   req.City,
		"key": cfg.Weather.API_KEY,
	}

	weather, err := service.FetchWeatherNow(cfg.Weather.API_URL, queryParams)
	if err != nil {
		errMsg := err.Error()

		switch {
		case strings.Contains(errMsg, "400"):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		case strings.Contains(errMsg, "404"):
			c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"description": "Successful operation - current weather forecast returned",
		"data":        weather,
	},
	)
}
