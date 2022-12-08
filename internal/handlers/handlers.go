package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/forms"
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
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.gohtml", &models.TemplateData{})
}

// Get request handler for make-reservation template
// yani tarayıcıda o template sayfası açıldığında bu fonksiyon çalışacak ve server tarafında empty Form objesi oluşturacak ve onu html'e gönderecek ve belirttiğimiz sayfayı render edecek..
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation // sayfa ilk yüklendiğinde emptyReservation gösterilecek
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// to handle the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	// bunu renderTemplate ile templateData verisine geçtik ve oradan html sayfalarımızda kullanabiliriz.
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone") // bunlardan herhangi biri empty ise hata mesajı ekliyoruz bu key değerlerine forms.go içinde.
	form.MinLength("first_name", 3, r)                         // first_name field en az 3 karakter olmalı
	form.IsEmail("email")                                      // email field gerçek bir email içermeli

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
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
