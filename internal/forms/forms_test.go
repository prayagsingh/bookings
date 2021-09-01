package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {

	// creating a request
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid form")
	}
}

func TestForm_Required(t *testing.T) {
	// creating a request
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	// test should fail if it shows form is valid
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "b")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("show does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	// creating a request
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")

	if has {
		t.Error("form shows had field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {

	// creating a request
	r := httptest.NewRequest("POST", "/whatever", nil)
	// initialize the form with empty values
	form := New(r.PostForm)

	form.MinLength("aaa", 13)
	// this should return error since we are checking the length of non-existence field
	if form.Valid() {
		t.Error("form shows min length for non-existence field")
	}

	// testing Get method in errors.go
	// case 1: where there is an error
	isError := form.Errors.Get("aaa")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	// creating a new form and passing postedData to the form and then recheck te min length
	postedData := url.Values{}
	postedData.Add("some_field", "some value")
	form = New(postedData)
	// len(some_field) must be less than 100
	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("show min length of 100 met when data is shorter")
	}

	// negative case.
	postedData = url.Values{}
	postedData.Add("some_field", "some value")
	form = New(postedData)

	form.MinLength("some_field", 1)
	if !form.Valid() {
		t.Error("show min length of 1 is not met when it is")
	}

	// case 2: where there is no error
	isError = form.Errors.Get("some_field")
	if isError != "" {
		t.Error("should not have an error, but did got one")
	}
}

func TestForm_IsEmail(t *testing.T) {

	// creating a request
	postedData := url.Values{}
	// initialize the form with empty values
	form := New(postedData)
	form.IsEmail("x")

	if form.Valid() {
		t.Error("form shows valid email for non-existence field")
	}

	postedData = url.Values{}
	postedData.Add("email", "abc@gmail.com")
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	// negative case. passing an invalid email address
	postedData = url.Values{}
	postedData.Add("email", "abc.com")
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got an valid resp for an invalid email address")
	}

}
