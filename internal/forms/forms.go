package forms

import (
	"fmt"
	"net/http"
	"net/url"
)

// Form creates a new form struct and embeds a url.values object
type Form struct {
	url.Values
	Errors errors
}

//Valid returns true if there are no errors
func(f *Form) Valid() bool {

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

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {

	x := r.Form.Get(field)
	fmt.Println("value of x is: ", x)
	if x == "" {

		f.Errors.Add(field, "This field can't be blank")
		return false
	}

	return true
}
