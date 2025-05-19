package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/weather", HandleGetWeather)

	router.GET("/subscribe", func(c *gin.Context) {
		c.HTML(http.StatusOK, "subscribe.html", nil)
	})
	router.POST("/subscribe", handleSubscribe)
	router.GET("/confirm/:token", handleConfirm)
	router.GET("/unsubscribe/:token", handleUnsubscribe)
}
