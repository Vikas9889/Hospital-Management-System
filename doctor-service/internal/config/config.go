package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	db := os.Getenv("DATABASE_URL")
	if db == "" {
		log.Fatal("DATABASE_URL not set")
	}

	return &Config{
		Port:        port,
		DatabaseURL: db,
	}
}
