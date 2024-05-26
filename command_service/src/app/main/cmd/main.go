package main

import (
	"context"
	"github.com/MercyClassic/go_saga/src/app/main/dependencies"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	ctx := context.Background()
	dependencies.Init(
		ctx,
		router,
		os.Getenv("db_uri"),
	)
	log.Println("Server started")
	log.Println(http.ListenAndServe("0.0.0.0:8001", router))
}
