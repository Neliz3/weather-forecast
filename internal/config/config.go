package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := getEnv("PORT", ":8080")

	return Config{
		Port: port,
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
