package handlers

import (
	"errors"
	"fmt"
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

func (m *Repository) Divide(w http.ResponseWriter, r *http.Request) {
	result, err := divideValues(100.0, 0)
	if err != nil {
		fmt.Fprintf(w, "Error occurred: %s\n", err)
		return
	}
	n, err := fmt.Fprintf(w, "%f / %f = %f", 100.0, 0.0, result) // verdiğimiz içeriği w'ya yazar. n kaç byte yazıldığıdır.
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Number of bytes written: %d\n", n)

}

func divideValues(x, y float32) (float32, error) {
	// Hata olması durumunda default 0 döndürüp, bir error dönüyoruz.
	if y == 0 {
		err := errors.New("0 ile bölme işlemi yapilamaz")
		return 0, err
	}
	// Hata olmazsa işlem sonucunu dönüyoruz ve err yerine nil değerini atıyoruz.
	result := x / y
	return result, nil
}
