package router

import (
	"database/sql"
	"net/http"

	"github.com/maxpn01/x-twitter-clone/auth"
	"github.com/maxpn01/x-twitter-clone/handler"
	"github.com/maxpn01/x-twitter-clone/repository"
)

func Router(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	authService := auth.NewAuthService(repository.NewUserRepository(db))
	authHandler := handler.NewAuthHandler(authService)

	mux.HandleFunc("GET /", handler.Home)
	mux.HandleFunc("POST /auth/signup", authHandler.Signup)

	return mux
}
