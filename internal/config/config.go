package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Weather struct {
	WEATHER_API_KEY string
	WEATHER_API_URL string
}

type Config struct {
	Port    string
	Weather Weather
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	PORT := getEnv("PORT", ":8080")
	WEATHER_API_KEY := getEnv("WEATHER_API_KEY", "")
	WEATHER_API_URL := getEnv("WEATHER_API_URL", "")

	return Config{
		Port:    PORT,
		Weather: Weather{WEATHER_API_KEY: WEATHER_API_KEY, WEATHER_API_URL: WEATHER_API_URL},
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
