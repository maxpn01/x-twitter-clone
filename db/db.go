package db

import (
	"database/sql"
	"log"
	"os"
)

func ConnectDB() *sql.DB {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
