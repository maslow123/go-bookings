package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maslow123/bookings/cmd/internal/config"
	"github.com/maslow123/bookings/cmd/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/generals-quarters", http.HandlerFunc(handlers.Repo.Generals))
	mux.Get("/majors-suite", http.HandlerFunc(handlers.Repo.Majors))
	mux.Get("/search-availability", http.HandlerFunc(handlers.Repo.Availability))
	mux.Post("/search-availability", http.HandlerFunc(handlers.Repo.PostAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(handlers.Repo.AvailabilityJSON))
	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.Reservation))
	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))

	fileServer := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/assets/*", http.StripPrefix("/assets", fileServer))

	return mux
}
