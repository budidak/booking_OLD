package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/driver"
	"github.com/budidak/booking/internal/handlers"
	"github.com/budidak/booking/internal/models"
	"github.com/budidak/booking/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig //mainde bunu oluşturup değerlerini veriyoruz ya da değiştiriyoruz, ama kullanılacak diğer dosyalarda da bunu pointer olarak *config.AppConfig olarak kullanıyoruz ki hepsi memoryde aynı şeyi göstersin.
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

// main function is the entry point to our program
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close() // main() sonlanana kadar database connection kesilmeyecek.

	// Starting server
	fmt.Printf("Starting application on port %s\n", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// main fonksiyonumuzu olabildiğince kısa yazmalıyız çünkü maini test etmeyeceğiz.
// Onun içerisindeki kodları ayrı bir fonksiyon olarak yazarak mainde çalıştırdık sadece.
func run() (*driver.DB, error) {
	// session ne tür bilgi tutacak (aynı post bilgisini farklı url'lerde kullanmak için)
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// app. olanlar configuration için
	app.InProduction = false

	// programımız için infoLog oluşturduk, ilk parametre terminale yazmamızı, ikincisi prefix, üçüncüsü de ekstra bilgi yani zaman vs.
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// programımız için error log oluşturduk.
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// session management
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=booking user=postgres password=postgres")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database.")
	// defer db.SQL.Close() // bunu burada run() içinde sonlandırırsak mainde kullanamayız, bu yüzden run() fonksiyonunun dönüş değeri olarak ekledik.

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		// return err
		return nil, err
	}
	app.TemplateCache = templateCache
	app.UseCache = false

	// Repository pattern
	// repo := handlers.NewRepo(&app) // returns App (so =>  repo.App = *config.AppConfig)
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo) // creates Repo variable = repo (so we can use Repo.App in handlers.go now)
	render.NewRenderer(&app)

	// return nil
	return db, nil // Buraya kadar çalışırsa zaten hata yoktur, nil dönebiliriz.
}
