package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/models"
	"github.com/budidak/booking/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplate string = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// ------------------------------------------------- //
	// handlers.go dosyamızdan önce main çalıştığı için main.go içindeki kodları buraya kopyaladım.
	// main.go içinde routes.go kullanıldığı için oradaki kodları da kopyaladım.
	// routes.go içinde middlewares.go kullanıldığı için oradaki kodları da kopyaladım.
	// çünkü handlers_test.go çalıştırılmadan önce bunlar istenecek, o yüzden onları burada hallettik.
	// ------------------------------------------------- //

	gob.Register(models.Reservation{})
	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	templateCache, err := CreateTestTemplateCache() // yeni oluşturduğumuz test fonksiyonu, en altta
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = templateCache
	app.UseCache = true // false yaparsak her request sayfa yeniden oluşturulur.

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf) // bu post request yapılan sayfalarda csrf token istediği için eğer o yoksa http response 200 yerine 400 dönebilir bad gateway, bunu test için comment yaptık, çünkü bunun testini zaten middleware yazarken yapmıştık.
	mux.Use(SessionLoad)
	// zaten handlers package içinde olduğumuz için handlers.Repo.Home yerine Repo.Home olarak çalıştırdık.
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	mux.Get("/contact", Repo.Contact)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.gotmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gotmpl", pathToTemplate))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gotmpl", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
