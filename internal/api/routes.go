package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	router.GET("/weather", handleGetWeather)
}
