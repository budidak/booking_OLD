package render

import (
	"net/http"
	"testing"

	"github.com/budidak/booking/internal/models"
)

// önce setup_test.go oluşturduk, sonra burada testlerimizi yazdık.

// AddDefaultData() fonksiyonumuzu çağırabilmek için gerekli parametreleri replike ettik.
// td kolayca oluşturduk ama r *http.Request kısmını http.NewRequest ile yaptık.
// Buradaki problem test çalıştırırken AddDefaultData() fonksiyondaki aynı sessiona sahip olamamamız
// Bu yüzden aşağıda getSession() yazdık, bu fonksiyon *http.Request, error dönebiliyor.
func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	// oluşturduğumuz session doğru çalışacak mı onu kontrol edelim.
	// session içerisine key:value olarak -> flash:12345 atadık.
	session.Put(r.Context(), "flash", "12345")

	result := AddDefaultData(&td, r)
	if result.Flash != "12345" {
		t.Error("flash value of 12345 not found in session")
	}
}

// *http.Request tipinde bir değiştirebilmek için yazdık bunu.
func getSession() (*http.Request, error) {
	// rastgele bir url için GET request oluşturduk ve başarılı olursa r değişkenine atadık.
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	// r değişkenimiz üzerinden bir Context oluşturduk ve bu contexte aynen aşağıdaki gibi session.Load() yapmak zorundayız.
	// daha sonra en son contextimizi tekrardan r değişkenimize atıyoruz.
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplate = "./../../templates" //package level variable olduğu için direkt değerini değiştritebildik.
	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = templateCache // app; package level variable direkt kullanabildik. (render.go içinde tanımlanmıştı)

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// *http.Request objemizi yazdığımız getSessions() ile modellemiştik,
	// w http.ResponseWriter objemiz için setup_test.go içerisinde bir şey yazdık. (myWriter)

	var ww myWriter

	err = Template(&ww, r, "home.gotmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to the browser")
	}

	err = Template(&ww, r, "non-existing.gotmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template that does not exist")
	}
}

func TestNewRenderer(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplate = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
