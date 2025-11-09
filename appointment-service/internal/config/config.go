package config

import (
	"log"
	"os"
)

type Config struct {
	Port             string
	DatabaseURL      string
	UserServiceURL   string
	DoctorServiceURL string
}

func Load() *Config {
	cfg := &Config{
		Port:             os.Getenv("PORT"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		UserServiceURL:   os.Getenv("USER_SERVICE_URL"),
		DoctorServiceURL: os.Getenv("DOCTOR_SERVICE_URL"),
	}

	if cfg.Port == "" || cfg.DatabaseURL == "" {
		log.Fatal("missing required environment variables")
	}

	return cfg
}
