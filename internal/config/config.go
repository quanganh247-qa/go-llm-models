package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	OllamaURL   string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		OllamaURL:   os.Getenv("OLLAMA_URL"),
		Port:        os.Getenv("PORT"),
	}
}
