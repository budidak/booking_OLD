package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/budidak/booking/internal/config"
)

var app *config.AppConfig

// sets config.AppConfig for this file. (mainde ayarlamalarımızı yaptıktan sonra en son orada bu fonksiyonu çağıracağız, buradaki dosyada da orada ayarlanan özellikler geçerli olacak.)
func NewHelpers(a *config.AppConfig) {
	app = a
}

// 2 farklı error olabilir: Client Error veya Server Error
// Bu fonksiyonlarda w http.ResponseWriter yazdık çünkü client'a bilgi vermek için kullanmayı düşündük.

func ClientError(w http.ResponseWriter, status int) {
	// yukarıda config ayarladığımız için burada app üzerinden config.AppConfig içindeki şeylere erişebiliyoruz. Örnek olarak InfoLog olarak yazdıralım.
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status) // feedback
}

func ServerError(w http.ResponseWriter, err error) {
	// oluşan hatayı detaylı bir şekilde yazdırmak isteyelim. Bunun için aşağıdaki gibi trace adında bir değişken oluşturduk.
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // feedback
}

// production aşamasında; hatalarımızı bir dosyaya yazdırmak isteyebiliriz. Kullanıcıları mail veya sms ile bilgilendirmek isteyebiliriz.

// developing aşamasında olduğumuz için şu an hatalarımızı basitçe terminale yazdırdık.

// diğer dosyalarımızdaki error kısımlarını buradaki yazdıklarımızla değiştirmek istiyoruz
// (handlers.go içindekileri falan) Örneğin PostReservation kısmında if err != nil {} vardı onun içindeki güncelledik. (log.Println() yerine -> helpers.ServerError())

/* codebase güncelledikten sonra (örneğin burada main.go içerisine aşağıdakileri eklemiştik)

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog)

İşimiz tamamlandığında önce testlerimizi çalıştırmalıyız. Root levele geldikten sonra tüm testleri çalıştırmak için şu komutu kullanabiliriz. > go test -v ./...
Burada bizde hata olan yerleri de gösterecek çünkü henüz testlerimizi update etmedik yeni kodlara göre.
Hataları gördük, handlers ve render klasörlerindeki setup_test.go dosyalarına mainde eklediğimiz bu kodu uygun şekilde eklemeliyiz. Örneğin handlers/ içinde app ismiyle kullandık ama renders/ içinde testApp olarak kullanmıştık kodumuzu bunu da ona göre düzenlemeliyiz.

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog)
*/
