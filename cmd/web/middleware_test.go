package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH) // replika ettiğimiz ServeHTTP metodu bir pointer üzerinden çağrılıyor diye pointer olarak geçtik.
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("Type is not http.Handler, but is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH) // replika ettiğimiz ServeHTTP metodu bir pointer üzerinden çağrılıyor diye pointer olarak geçtik.
	// Fonksiyonu çalıştırdıkdan sonra dönen sonuç h. Yeni bir v değişkeni tanımlıyoruz h ile aynı tipte, eğer bu tip http.Handler ise sorun yok. Değilse hata mesajı basıyoruz.
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("Type is not http.Handler, but is %T", v)
	}

}
