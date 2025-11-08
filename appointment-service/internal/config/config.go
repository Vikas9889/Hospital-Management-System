package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port           string
    DatabaseURL    string
    UserServiceURL string
}

func Load() *Config {
    _ = godotenv.Load()

    cfg := &Config{
        Port:           getEnv("PORT", "8082"),
        DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres-appointments:5432/appointments_db?sslmode=disable"),
        UserServiceURL: getEnv("USER_SERVICE_URL", "http://user-service:8081"),
    }

    log.Println("Config loaded:", cfg)
    return cfg
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
