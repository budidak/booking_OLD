package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// creates a custom Form struct and it embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}), // errors{} ile aynı şey sadece empty errors objesi oluşturduk.
	}
}

// Valid returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// form ile alınan bilginin name değerinin, serverda karşılığı var mı?
func (f *Form) Has(field string) bool {
	// x := r.Form.Get(field)  // bu varken parametrede r *http.Request vardı ama artık kaldırdık gerek yok.
	x := f.Get(field) // böyle yapmalıyız çünkü yukarıdaki gibi yapınca request içinden forma erişmeye çalışıyorsa testten kalıyor. (forms_test.go)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank.")
		return false
	}
	return true
}

// ... kullandığımız zaman variadic function olarak isimlendiriliyor. *args
// This checks if the fields are empty
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		// trims the extra spaces
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

// checks the minimum length of the fields
func (f *Form) MinLength(field string, length int) bool {
	// x := r.Form.Get(field)
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// checks if the email address is valid
// https://github.com/asaskevich/govalidator
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address.")
	}
}
