package models

import "github.com/budidak/booking/internal/forms"

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string // for post request csrf token
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form // client-server arasında kullanacağımız form bilgileri için.
}
