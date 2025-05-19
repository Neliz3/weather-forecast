package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	PORT := getEnv("PORT", ":8080")
	BaseURL := getEnv("BASE_URL", "http://localhost:8080/api")

	WEATHER_API_KEY := getEnv("WEATHER_API_KEY", "")
	WEATHER_API_URL := getEnv("WEATHER_API_URL", "")

	EMAIL_FROM := getEnv("EMAIL_FROM", "")
	EMAIL_API_KEY := getEnv("EMAIL_API_KEY", "")
	SECRET_KEY_JWT := getEnv("SECRET_KEY_JWT", "")

	DB_HOST := getEnv("DB_HOST", "localhost")
	DB_PORT := getEnv("DB_PORT", "5432")
	DB_USER := getEnv("POSTGRES_USER", "")
	DB_PASSWORD := getEnv("POSTGRES_PASSWORD", "")
	DB_NAME := getEnv("POSTGRES_NAME", "")

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

func getEnv[T any](key string, defaultVal T) T {
	valStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	switch any(defaultVal).(type) {
	case int:
		if val, err := strconv.Atoi(valStr); err == nil {
			return any(val).(T)
		}
	case string:
		return any(valStr).(T)
	}

	return defaultVal
}
