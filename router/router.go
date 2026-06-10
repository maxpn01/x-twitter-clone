package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/maxpn01/x-twitter-clone/user"
	"github.com/maxpn01/x-twitter-clone/user/auth"
)

func Router(db *sql.DB) http.Handler {
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

	mux.HandleFunc("GET /", healthHandler)
	mux.HandleFunc("GET /api", healthHandler)

	mux.HandleFunc("POST /api/signup", userHandler.Signup)
	mux.HandleFunc("POST /api/signin", userHandler.Signin)
	mux.HandleFunc("POST /api/signout", userHandler.Signout)

	return recoverMiddleware(mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v\n%s", err, debug.Stack())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
