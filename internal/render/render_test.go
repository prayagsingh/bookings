package render

import (
	"net/http"
	"testing"

	"github.com/prayagsingh/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {

	var td models.TemplateData
	// creating session
	r, err := getSession()
	if err != nil {
		t.Log(err)
		t.Error(err)
	}

	// testing session
	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 is not found in session")
	}
}

// creating a func so that our request has a session data
func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	// we need to create a context else if we return r here then it will result in no session data in context error
	ctx := r.Context()
	// putting session data into the context
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	// putting context back into the request
	r = r.WithContext(ctx)

	return r, nil
}
