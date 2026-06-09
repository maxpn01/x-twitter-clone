package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/maxpn01/x-twitter-clone/db"
	"github.com/maxpn01/x-twitter-clone/router"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	godotenv.Load()

	// db connection

	databaseURL := os.Getenv("DATABASE_URL")

	db := db.ConnectDB(databaseURL)
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to postgres database")

	// migrations

	m, err := migrate.New("file://migrations", strings.Replace(databaseURL, "postgres://", "pgx5://", 1))
	if err != nil {
		log.Fatal("migration init error: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration error: ", err)
	}

	log.Println("migrations applied")

	// http server

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	router := router.Router(db)

	log.Println(("server listening on http://localhost" + addr))
	log.Fatal(http.ListenAndServe(addr, router))
}
