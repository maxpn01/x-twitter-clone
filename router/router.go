package router

import (
	"database/sql"
	"fmt"
	"net/http"
)

func Router(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello from go")
	})

	return mux
}
