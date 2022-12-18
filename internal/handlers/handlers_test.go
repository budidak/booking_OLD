package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// form verisinde neyi post ettiğimizi key:value şeklinde tutacak
type postData struct {
	key   string
	value string
}

// post ve get requestlerimizle eriştiğimiz routes için test table oluşturuyoruz.
var theTests = []struct {
	name               string     // oluşturduğumuz her bir table elemanına isim verdik
	url                string     // eriştiğimiz url
	method             string     // url'ye hangi methodla eriştiğimiz
	params             []postData // post method ile url'ye geçtiğimiz veriler
	expectedStatusCode int        // sorun yoksa 200 OK, sayfa bulunamazsa 404 Not Found, veya redirect varsa 300lü kodlar görürüz
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "01-01-2020"},
		{key: "end", value: "02-01-2020"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "01-01-2020"},
		{key: "end", value: "02-01-2020"},
	}, http.StatusOK},
	{"post-make-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Burak"},
		{key: "last_name", value: "Yesilyurt"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			// POST
			values := url.Values{} // testServer'imize post edilecek veriler için boş bir yapı oluşturduk, parametreleri key value olarak bunun içine atacağız.
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
