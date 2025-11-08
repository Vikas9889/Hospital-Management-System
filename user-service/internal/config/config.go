package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DatabaseURL string
}

func Load() *Config {
    _ = godotenv.Load()
    port := os.Getenv("PORT")
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL not set")
    }
    if port == "" {
        port = "8081"
    }
    return &Config{Port: port, DatabaseURL: dbURL}
}
