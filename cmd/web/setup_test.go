package main

import (
	"net/http"
	"os"
	"testing"
)

// Buradaki şeyler, testler çalıştırılmadan önce çalıştırılacak.
func TestMain(m *testing.M) {

	os.Exit(m.Run()) // m.Run() testleri çalıştırır ve daha sonra programdan çıkarız.
}

// middleware.go içindeki fonksiyonlarımız next http.Handler tipinde parametre alıp http.Handler tipinde return ediyordu.
// http.Handler aslında ServeHTTP metodunu implement eden bir interface olduğu için onu replika etmek için böyle bir kod yazdık.
type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
