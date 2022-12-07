package handlers

import (
	"net/http"

	"github.com/budidak/booking/pkg/config"
	"github.com/budidak/booking/pkg/models"
	"github.com/budidak/booking/pkg/render"
)

// config.go düzenlemesinden sonra yazdık bu kısmı da.
// REPOSITORY PATTERN: Allows us to swap components out of our application with a minimal changes required to the code base.

var Repo *Repository // repository used by the handlers (tüm handler metotlarımıza bunun üzerinden erişebileceğiz)

type Repository struct {
	App *config.AppConfig
}

// creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Repository pattern yazdıktan sonra bunları receiver function haline getiriyoruz.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Home ziyaret edildiğinde requestten ip adresini alıp sessiona kaydediyoruz örneğin.
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP) // key:value şeklinde session'a bilgi ekliyoruz.
	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	// Home functionda bilgiyi yükledik, burada sessiondan bilgiyi alabiliyoruz.
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.gohtml", &models.TemplateData{})
}
