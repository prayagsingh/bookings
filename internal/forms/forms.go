package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Form creates a new form struct and embeds a url.values object
type Form struct {
	url.Values
	Errors errors
}

//Valid returns true if there are no errors
func (f *Form) Valid() bool {

	// returns true if error is empty
	return len(f.Errors) == 0
}

// New initializes the form struct
func New(data url.Values) *Form {

	return &Form{
		data,
		// can't put the errors like (map[string][]string) because we are declaring it to be empty
		// so we have to put {} here.
		errors(map[string][]string{}),
	}
}

//Required is variadic func which take n number of field args and checks for these field
func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can't be blank")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {

	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength checks for string min length that the client is puting in the name fields
func (f *Form) MinLength(field string, length int, r *http.Request) bool {

	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}
