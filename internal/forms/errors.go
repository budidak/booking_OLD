package forms

// created custom type of map
type errors map[string][]string

// adds an error message for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// returns the first error message (bir field için birden fazla hata olabilir o yüzden her seferinde arraydaki ilk elemanı gösterdim.)
func (e errors) Get(field string) string {
	es := e[field]
	if len(e) == 0 {
		return ""
	}
	return es[0]
}
