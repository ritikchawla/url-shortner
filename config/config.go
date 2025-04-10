package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	RedisHost  string
	RedisPort  string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// It's ok if .env doesn't exist, we'll use environment variables
		fmt.Println("Warning: .env file not found")
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "urluser"),
		DBPassword: getEnv("DB_PASSWORD", "urlpass"),
		DBName:     getEnv("DB_NAME", "urlshortener"),
		DBPort:     getEnv("DB_PORT", "5432"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
