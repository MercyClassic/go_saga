package api

import (
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/MercyClassic/go_saga/src/app/presentators/api/v1"
	"github.com/go-chi/chi/v5"

	"log"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Println(err)
	}
}

func IncludeRouters(r chi.Router, pool client.Client) {
	r.Get("/ping", ping)
	v1.IncludeCommandRouter(r, pool)
}
