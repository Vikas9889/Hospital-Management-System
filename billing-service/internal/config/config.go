package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port                string
    DatabaseURL         string
    AppointmentServiceURL string
    UserServiceURL      string
    NotificationURL     string
}

func Load() *Config {
    _ = godotenv.Load()
    port := os.Getenv("PORT")
    dbURL := os.Getenv("DATABASE_URL")
    appt := os.Getenv("APPOINTMENT_SERVICE_URL")
    user := os.Getenv("USER_SERVICE_URL")
    notif := os.Getenv("NOTIFICATION_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL not set")
    }
    if appt == "" {
        log.Fatal("APPOINTMENT_SERVICE_URL not set")
    }
    if user == "" {
        log.Fatal("USER_SERVICE_URL not set")
    }
    if notif == "" {
        log.Fatal("NOTIFICATION_URL not set")
    }
    if port == "" {
        port = "8083"
    }
    return &Config{Port: port, DatabaseURL: dbURL, AppointmentServiceURL: appt, UserServiceURL: user, NotificationURL: notif}
}
