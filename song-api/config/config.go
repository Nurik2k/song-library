package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DatabaseURL    string
	ExternalAPIURL string
	Port           string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Ошибка загрузки .env файла: %v", err)
	}

	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/song_library"),
		ExternalAPIURL: getEnv("EXTERNAL_API_URL", "http://localhost:5000"),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
