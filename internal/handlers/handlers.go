package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/driver"
	"github.com/budidak/booking/internal/forms"
	"github.com/budidak/booking/internal/helpers"
	"github.com/budidak/booking/internal/models"
	"github.com/budidak/booking/internal/render"
	"github.com/budidak/booking/internal/repository"
	"github.com/budidak/booking/internal/repository/dbrepo"
)

// config.go düzenlemesinden sonra yazdık bu kısmı da.
// REPOSITORY PATTERN: Allows us to swap components out of our application with a minimal changes required to the code base.

var Repo *Repository // repository used by the handlers (tüm handler metotlarımıza bunun üzerinden erişebileceğiz)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo // database connection
}

// creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a), // database connection
	}
}

// sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Repository pattern yazdıktan sonra bunları receiver function haline getiriyoruz.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.gotmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.gotmpl", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.gotmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.gotmpl", &models.TemplateData{})
}

// Get request handler for make-reservation template
// yani tarayıcıda o template sayfası açıldığında bu fonksiyon çalışacak ve server tarafında empty Form objesi oluşturacak ve onu html'e gönderecek ve belirttiğimiz sayfayı render edecek..
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation // sayfa ilk yüklendiğinde emptyReservation gösterilecek
	render.Template(w, r, "make-reservation.gotmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// to handle the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Bunları form verisinde kullanıyoruz, databasede, rezervasyonla ilgili bilgiler.
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	// 01-01-2022 gibi dönecek bunu handle etmeliyiz önce. GO date&time formatlarken garip bir yaklaşım sergiler.
	// https://www.pauladamsmith.com/blog/2011/05/go_time.html
	timeLayout := "2006-01-02"
	startDate, err := time.Parse(timeLayout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(timeLayout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	// bunu Template() ile templateData verisine geçtik ve oradan html sayfalarımızda kullanabiliriz.
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email") // bunlardan herhangi biri empty ise hata mesajı ekliyoruz bu key değerlerine forms.go içinde.
	form.MinLength("first_name", 3)                   // first_name field en az 3 karakter olmalı
	form.IsEmail("email")                             // email field gerçek bir email içermeli

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.gotmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// write to database
	err = m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
	}

	// buradan yollanan post bilgilerini, başka bir sayfada almak için sessions kullanabiliriz.
	m.App.Session.Put(r.Context(), "reservation", reservation) // direkt bunu yapamayız, mainde sessiona ne tür bilgi koyacağımızı belirtmeliyiz.
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Post işleminden sonra rezervasyon bilgilerini göstereceğimiz sayfanın handler fonksiyonu
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) // dönecek bilginin tipini de en sonda .() ile belirttik.
	if !ok {
		// log.Println("Cannot get item from session.")
		m.App.ErrorLog.Println("Cannot get item from session.")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // hata olursa direkt anasayfaya yönlendiriyoruz.
		return
	}
	m.App.Session.Remove(r.Context(), "reservation") // bilgiyi aldıktan sonra sessiondan kaldırmalıyız böyle.
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.gotmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.gotmpl", &models.TemplateData{})
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
		helpers.ServerError(w, err)
		return
	}
	log.Println(string(out)) // Dönen bilgiyi terminalde görmek için.
	// tarayıcıya ne tür bir response döndüğümüzü söylemek için header yazmalıyız.
	w.Header().Set("Content-type", "application/json")
	w.Write(out) // Write() sends our response to response writer.
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.gotmpl", &models.TemplateData{})
}
