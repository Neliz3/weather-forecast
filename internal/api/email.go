package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"weather-forecast/internal/config"
	"weather-forecast/internal/service"
)

type SubscribeRequest struct {
	Email     string `form:"email" binding:"required,email"`
	City      string `form:"city" binding:"required"`
	Frequency string `form:"frequency" binding:"required,oneof=hourly daily"`
}

type UnsubscribeRequest struct {
	Email string `form:"email" binding:"required,email"`
}

func handleSubscribe(c *gin.Context) {
	var req SubscribeRequest
	cfg := config.Load()

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for existing subscription in DB

	// Generate token with email, city, frequency
	token, err := service.GenerateConfirmationToken(req.Email, req.City, req.Frequency, cfg.Email.SECRET_KEY_JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	if err := service.SendConfirmationEmail(cfg.Email.EmailFrom, req.Email, token, cfg.Email.API_KEY, cfg.BaseURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation email sent"})
}

func handleConfirm(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}

	cfg := config.Load()
	email, city, frequency, err := service.ValidateConfirmationToken(token, cfg.Email.SECRET_KEY_JWT)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Store the subscription - you’ll eventually insert into a DB.
	// For now, we’ll simulate this with a print/log.
	fmt.Printf("Confirmed subscription: %s - %s (%s)\n", email, city, frequency)

	// In production: check if already subscribed before saving

	c.JSON(http.StatusOK, gin.H{
		"message":   "Subscription confirmed!",
		"email":     email,
		"city":      city,
		"frequency": frequency,
	})
}

func handleUnsubscribe(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}

	cfg := config.Load()
	email, _, _, err := service.ValidateConfirmationToken(token, cfg.Email.SECRET_KEY_JWT)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Simulate unsubscribing (e.g., DB delete or flag update)
	fmt.Printf("Unsubscribed: %s\n", email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Unsubscribed successfully.",
		"email":   email,
	})
}
