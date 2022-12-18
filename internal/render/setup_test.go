package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// setup için main.go içerisinden bazı kısımları aldık render.go ilgilendiren.
	gob.Register(models.Reservation{})

	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// http.ResponseWriter incelediğimizde onun aslında bir interface{} olduğunu ve
// Header() Header
// Write([]byte) (int, error)
// WriteHeader(statusCode int)
// metotlarını implement etme koşulu olduğunu görüyoruz. Yani kendimiz herhangi bir tipte obje oluşturup (interface olduğu için struct), ona bu metotları implement ettirirsek aslında o oluşturduğumuz objeye de http.ResponseWriter olarak davranabiliriz. Bunu aşağıda yapalım.

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func (tw *myWriter) WriteHeader(statusCode int) {

}
