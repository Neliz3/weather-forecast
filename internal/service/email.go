package service

import (
	"fmt"
	"time"
	"weather-forecast/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func GenerateConfirmationToken(email, city, frequency, secret string) (string, error) {
	claims := jwt.MapClaims{
		"email":     email,
		"city":      city,
		"frequency": frequency,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateConfirmationToken(tokenStr, secret string) (string, string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", "", "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", fmt.Errorf("invalid claims")
	}

	email, ok1 := claims["email"].(string)
	city, ok2 := claims["city"].(string)
	frequency, ok3 := claims["frequency"].(string)

	if !ok1 || !ok2 || !ok3 {
		return "", "", "", fmt.Errorf("missing claims")
	}

	return email, city, frequency, nil
}

func SendConfirmationEmail(fromEmail, toEmail, token, api_key, baseURL string) error {
	from := mail.NewEmail("Weather Forecast Service", fromEmail)
	to := mail.NewEmail("", toEmail)
	subject := "Welcome! Please confirm your subscription"
	url := fmt.Sprintf("%s/confirm/%s", baseURL, token)

	plain := fmt.Sprintf("Click the link to confirm: %s", url)
	html := fmt.Sprintf("<p>Click <a href=\"%s\">here</a> to confirm your subscription.</p>", url)

	message := mail.NewSingleEmail(from, subject, to, plain, html)
	client := sendgrid.NewSendClient(api_key)
	resp, err := client.Send(message)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("SendGrid error: %s", resp.Body)
	}
	fmt.Printf("Email sent. Status: %d\nBody: %s\n", resp.StatusCode, resp.Body)
	return nil
}

func SendWeatherUpdateEmail(fromEmail, toEmail, city, frequency string, weather map[string]any, cfg *config.Config) error {
	from := mail.NewEmail("Weather Forecast Service", fromEmail)
	to := mail.NewEmail("", toEmail)
	subject := "Your Weather Update"

	// Retrieve unsubscribed token
	token, err := GenerateConfirmationToken(toEmail, city, frequency, cfg.Email.SECRET_KEY_JWT)
	if err != nil {
		return fmt.Errorf("token generation failed: %s", err)
	}

	temperature := weather["temperature"]
	humidity := weather["humidity"]
	description := weather["description"]

	unsubscribeURL := fmt.Sprintf("%s/api/unsubscribe/%s", cfg.BaseURL, token)

	plainTextContent := fmt.Sprintf(
		"Weather update:\nTemperature: %v\nHumidity: %v\nDescription: %v\n\nTo unsubscribe, click: %s",
		temperature, humidity, description, unsubscribeURL,
	)
	htmlContent := fmt.Sprintf(
		`<p>Weather update:</p>
		<ul>
			<li>Temperature: %v</li>
			<li>Humidity: %v</li>
			<li>Description: %v</li>
		</ul>
		<p><a href="%s" style="color: red;">Unsubscribe from these emails</a></p>`,
		temperature, humidity, description, unsubscribeURL,
	)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(cfg.Email.API_KEY)
	resp, err := client.Send(message)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("SendGrid error: %s", resp.Body)
	}

	return nil
}
