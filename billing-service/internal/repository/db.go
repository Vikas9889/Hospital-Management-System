package repository

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

func ConnectDB(url string) *sql.DB {
    db, err := sql.Open("postgres", url)
    if err != nil {
        log.Fatalf("Failed to connect DB: %v", err)
    }
    if err := db.Ping(); err != nil {
        log.Fatalf("DB unreachable: %v", err)
    }
    return db
}
