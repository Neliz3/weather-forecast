package config

import (
	"os"
)

type Weather struct {
	API_KEY string
	API_URL string
}

type Email struct {
	EmailFrom      string
	API_KEY        string
	SECRET_KEY_JWT string
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	Port    string
	BaseURL string
	Weather Weather
	Email   Email
	DB      DB
}

func Load() *Config {
	PORT := os.Getenv("PORT")
	BaseURL := os.Getenv("BASE_URL")

	WEATHER_API_KEY := os.Getenv("WEATHER_API_KEY")
	WEATHER_API_URL := os.Getenv("WEATHER_API_URL")

	EMAIL_FROM := os.Getenv("EMAIL_FROM")
	EMAIL_API_KEY := os.Getenv("EMAIL_API_KEY")
	SECRET_KEY_JWT := os.Getenv("SECRET_KEY_JWT")

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("POSTGRES_USER")
	DB_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	DB_NAME := os.Getenv("POSTGRES_DB")

	return &Config{
		Port:    PORT,
		BaseURL: BaseURL,
		Weather: Weather{
			API_KEY: WEATHER_API_KEY,
			API_URL: WEATHER_API_URL},
		Email: Email{
			EmailFrom:      EMAIL_FROM,
			API_KEY:        EMAIL_API_KEY,
			SECRET_KEY_JWT: SECRET_KEY_JWT,
		},
		DB: DB{
			Host:     DB_HOST,
			Port:     DB_PORT,
			User:     DB_USER,
			Password: DB_PASSWORD,
			DBName:   DB_NAME,
		},
	}
}
