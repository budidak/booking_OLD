package main

import (
	"net/http"

	"github.com/budidak/booking/pkg/config"
	"github.com/budidak/booking/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// go get -u github.com/go-chi/chi

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	// middleware allows you to process a request as it comes into your Web applicaton and performs some action on it.
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	// Projemizdeki statik dosyaları (image, css, js göstermek için kullanıyoruz.)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // static/ klasöründeki her şeyi seçer ve onlardan /static prefixini kaldırır.
	return mux

}
