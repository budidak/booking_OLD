package main

import (
	"net/http"

	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// go get -u github.com/go-chi/chi

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	// middleware allows you to process a request as it comes into your Web applicaton and performs some action on it.
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf) // POST requestlerde CSRF protection için bunu hep kullan. Uygun bir CSRF token üretmeyen post requestler engellenir.
	mux.Use(SessionLoad)

	// routing
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability) // form verisini post request ile yollamak için.
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/contact", handlers.Repo.Contact)

	// Projemizdeki statik dosyaları (image, css, js göstermek için kullanıyoruz.)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // static/ klasöründeki her şeyi seçer ve onlardan /static prefixini kaldırır.
	return mux

}
