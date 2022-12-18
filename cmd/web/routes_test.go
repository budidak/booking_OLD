package main

import (
	"testing"

	"github.com/budidak/booking/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app) // routes fonksiyonunu alması gerektiği parametre ile çağırdık ve döndüğü değeri aldık.

	// dönen değişken dönmesi gereken tipte mi diye kontrol ettik, değilse hata verdik.
	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("Type is not not *chi.Mux, type is %T", v)
	}
}
