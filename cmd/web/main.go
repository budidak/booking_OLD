package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/handlers"
	"github.com/budidak/booking/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig //mainde bunu oluşturup değerlerini veriyoruz ya da değiştiriyoruz, ama kullanılacak diğer dosyalarda da bunu pointer olarak *config.AppConfig olarak kullanıyoruz ki hepsi memoryde aynı şeyi göstersin.
var session *scs.SessionManager

func main() {
	// app. olanlar configuration için
	app.InProduction = false

	// session management
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = templateCache
	app.UseCache = false

	// Repository pattern
	repo := handlers.NewRepo(&app) // returns App (so =>  repo.App = *config.AppConfig)
	handlers.NewHandlers(repo)     // creates Repo variable = repo (so we can use Repo.App in handlers.go now)
	render.NewTemplates(&app)

	// Starting server
	fmt.Printf("Starting application on port %s\n", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	log.Fatal(srv.ListenAndServe())
}
