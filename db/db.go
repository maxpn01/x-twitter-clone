package db

import (
	"database/sql"
	"log"
)

func ConnectDB(databaseURL string) *sql.DB {
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
