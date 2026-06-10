package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxpn01/x-twitter-clone/user"
	"github.com/maxpn01/x-twitter-clone/user/auth"
)

func Router(db *sql.DB) *http.ServeMux {
	repo := &user.PostgresUserRepository{DB: db}
	jwtService, err := auth.NewJWTService(os.Getenv("JWT_SECRET"), 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	userService, err := user.NewUserService(repo, jwtService)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	userHandler := user.NewUserHandler(userService)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "the server is alive ")
	})

	mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "the api is running and healthy")
	})

	mux.HandleFunc("POST /api/signup", userHandler.Signup)

	return mux
}
