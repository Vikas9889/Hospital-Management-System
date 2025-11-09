package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// ConnectDB opens a Postgres connection and pings it. Pass a full Postgres URL.
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

type PrescriptionRepository struct{ db *sql.DB }

func NewPrescriptionRepository(db *sql.DB) *PrescriptionRepository {
	return &PrescriptionRepository{db: db}
}

// Add methods: Create, Get, List
