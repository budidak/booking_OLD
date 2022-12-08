package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// POST requestler yapabilmek için unique CSRF token oluşturacak sayfada.
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	// get the template cache from the app config if UseCache true, else read it from disk
	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}
}

// Dosya isimlerini tek tek manual yazmak yerine, Glob() ile dizin içinde arattık ve otomatikleştirdik.
// Yine cache map şeklinde key değeri string (muhtemel sayfa adı home.page.gohtml), value değeri de *template.Template tipinde olacak yani ilgili sayfanın taslağı tutulacak.
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{} // myCache := make(map[string]*template.Template)

	// get all files name *.page.gohtml from ./templates folder
	pages, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		return myCache, err
	}

	// loop through the pages, dizinde bulduğumuz dosya isimlerinin tutulduğu slice üzerinde döngü kurarak hepsinden ayrı ayrı Template oluşturuyoruz bunu da ParseFiles() ile yapıyoruz.
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// base layout var mı diye kontrol ediyoruz (çünkü diğer .gohtml sayfaları import edebilir.)
		matches, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return myCache, err
		}

		// eğer base layout varsa, page.gohtml sayfalarında gereklilik olacağı için ParseGlob() ile parse ediyoruz.
		// ilk versiyonlarda ("./templates/"+tmpl, "./templates/base.layout.gohtml") yaptığımız şeyi yaptık yani.
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return myCache, err
			}
		}
		// en son hiçbir hata yoksa elde edilen sayfayı cache değişkenimize ekliyoruz.
		myCache[name] = ts
	}
	return myCache, nil
}
