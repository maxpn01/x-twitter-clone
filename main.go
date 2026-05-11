package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/maxpn01/x-twitter-clone/db"
	"github.com/maxpn01/x-twitter-clone/router"
)

func main() {
	godotenv.Load()

	db := db.ConnectDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to postgres database")

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
