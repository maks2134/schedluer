package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	BSUIRAPI BSUIRAPIConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type BSUIRAPIConfig struct {
	BaseURL string
	Timeout int // в секундах
}

type LoggerConfig struct {
	Level string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", ""),
			Database: getEnv("MONGODB_DATABASE", "schedluer"),
		},
		BSUIRAPI: BSUIRAPIConfig{
			BaseURL: getEnv("BSUIR_API_BASE_URL", "https://iis.bsuir.by/api/v1"),
			Timeout: 30,
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	if config.MongoDB.URI == "" {
		return nil, fmt.Errorf("MONGODB_URI is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
