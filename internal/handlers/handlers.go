package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/models"
	"github.com/budidak/booking/internal/render"
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
	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	// Home functionda bilgiyi yükledik, burada sessiondan bilgiyi alabiliyoruz.
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// to handle POST request
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start") // Form bilgisini requestten bu şekilde alıyoruz, form içindeki elementin (input name ismidir)
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Posted to search availability : Start date is %s and End date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out)) // Dönen bilgiyi terminalde görmek için.
	// tarayıcıya ne tür bir response döndüğümüzü söylemek için header yazmalıyız.
	w.Header().Set("Content-type", "application/json")
	w.Write(out) // Write() sends our response to response writer.

}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}
