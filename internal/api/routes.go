package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/weather", handleGetWeather)
	router.POST("/subscribe", handleSubscribe)
	router.GET("/confirm/:token", handleConfirm)
	router.GET("/unsubscribe/:token", handleUnsubscribe)
}
