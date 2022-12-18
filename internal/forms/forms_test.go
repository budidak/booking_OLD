package forms

import (
	"net/url"
	"testing"
)

// bunlar receiver functions olduğu için isimlerini TestReceiverName_MethodName() şeklinde yazdık.

func TestForm_Valid(t *testing.T) {
	// Valid() metodu bir f *Form receiver üzerinden çağrıldığı için öyle bir şeyi replike etmeliyiz önce.
	// r := httptest.NewRequest("POST", "/whatever", nil)
	// form := New(r.PostForm)
	// veya boş bir url.Values{} oluşturmalı ve onun üzerinden New() ile form oluşturmalıyız.
	postedData := url.Values{}
	form := New(postedData)
	isValid := form.Valid()
	if !isValid {
		t.Error("invalid döndü, valid dönmesi bekleniyordu")
	}
}

// Test yazdığımız fonksiyonlarda hem valid durumu hem de invalid durumu değerlendirmeliyiz.
func TestForm_Required(t *testing.T) {
	// Required() metodu bir f *Form receiver üzerinden çağrıldığı için öyle bir şeyi replike etmeliyiz önce.
	postedData := url.Values{}
	form := New(postedData)
	form.Required("a", "b", "c") // variadic olarak yazdığımız için sınırsız parametre alabilir, bu fields gerekli dedik.
	if form.Valid() {            // eğer fields eklenmemişken form Valid() fonksiyonunu sağlarsa o zaman hata var demektir.
		t.Error("form shows valid when required fields missing")
	}

	// post edilecek verilerimizi ekliyoruz.
	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	// Eğer forma r üzerinden yeni bir veri ekleyecek olsaydık New(r.PostForm) kullanırdık, ama sade bir kod için r kullanmıyoruz.
	// r = httptest.NewRequest("POST", "/whatever", nil)
	// r.PostForm = postedData
	// form = New(r.PostForm)
	form = New(postedData)
	form.Required("a", "b", "c") // örneğin bu alanlar gerekli, ve bizim formumuzda da key değeri olarak bu fields eklendi yukarıda.
	if !form.Valid() {           // fields oluşturup postedDatamıza ekledik, eğer formumuz bu fields varken Valid() fonksiyonunu sağlayamazsa o zaman hata var demektir.
		t.Error("form does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	// r := httptest.NewRequest("POST", "whatever", nil)
	// form := New(r.PostForm)
	// has := form.Has("a", r) // has form field a?
	postedData := url.Values{}
	form := New(postedData)
	has := form.Has("a")
	if has {
		t.Error("formda bu field olmamasi gerekirken var gözüküyor")
	}

	postedData.Add("a", "a") // forma eklenecek test datasına yeni bir key:value çifti ekledik.
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("formda olmasi beklenen field, yok olarak gözüküyor")
	}
}

func TestForm_MinLength(t *testing.T) {
	// r := httptest.NewRequest("POST", "whatever", nil)
	// form := New(r.PostForm)
	postedData := url.Values{}
	form := New(postedData)
	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows minlength for non existing field")
	}

	postedData.Add("some_field", "some value")
	form = New(postedData)
	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("form shows minlenght of 100 met, when data is shorter")
	}

	postedData = url.Values{} // reinitializes as empty object
	postedData.Add("another_field", "abc123")
	form = New(postedData)
	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("form shows minlength of 1 is not met, but it is")
	}
}

func TestForm_Email(t *testing.T) {
	// r := httptest.NewRequest("POST", "whatever", nil)
	// form := New(r.PostForm)
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid for non-existing field")
	}

	postedData.Add("email", "me@here.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("form shows invalid email for valid email")
	}

	postedData = url.Values{}
	postedData.Add("email", "x")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid for invalid email address")
	}
}
