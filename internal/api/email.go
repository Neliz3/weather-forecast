package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"weather-forecast/internal/config"
	"weather-forecast/internal/db"
	"weather-forecast/internal/model"
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
	db_connect, _ := db.Connect()
	defer func() {
		if err := db_connect.Close(); err != nil {
			log.Printf("Failed to close DB connection: %v", err)
		}
	}()

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if already subscribed logic (fetch from DB)
	existing, err := db.GetSubscriptionByEmailCity(db_connect, req.Email, req.City)
	if err == nil && existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already subscribed"})
		return
	}

	// Insert subscription with confirmed = false
	err = db.InsertSubscription(db_connect, model.Subscription{
		Email:     req.Email,
		City:      req.City,
		Frequency: req.Frequency,
		Confirmed: false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save subscription: " + err.Error()})
		return
	}

	// Generate token with email, city, frequency
	token, err := service.GenerateConfirmationToken(req.Email, req.City, req.Frequency, cfg.Email.SECRET_KEY_JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	if err := service.SendConfirmationEmail(cfg.Email.EmailFrom, req.Email, token, cfg.Email.API_KEY, cfg.BaseURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})
}

func handleConfirm(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	cfg := config.Load()
	db_connect, _ := db.Connect()
	defer func() {
		if err := db_connect.Close(); err != nil {
			log.Printf("Failed to close DB connection: %v", err)
		}
	}()

	email, _, _, err := service.ValidateConfirmationToken(token, cfg.Email.SECRET_KEY_JWT)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	// Update DB to mark subscription as confirmed
	err = db.MarkSubscriptionConfirmed(db_connect, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription confirmed successfully",
	})
}

func handleUnsubscribe(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	cfg := config.Load()
	db_connect, _ := db.Connect()
	defer func() {
		if err := db_connect.Close(); err != nil {
			log.Printf("Failed to close DB connection: %v", err)
		}
	}()

	email, _, _, err := service.ValidateConfirmationToken(token, cfg.Email.SECRET_KEY_JWT)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	// Remove subscription from DB
	err = db.DeleteSubscription(db_connect, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Unsubscribed successfully",
	})
}
