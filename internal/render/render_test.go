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

func TestNewRender(t *testing.T) {

	NewRenderer(app)

}

func TestRenderTemplate(t *testing.T) {

	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	// adding tc to app variable
	app.TemplateCache = tc
	// get request using getSession
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter
	// please make sure you have all the required inputs such as http.ResponseWriter, Request etc
	// we you can access them directly here then need to create in setup_test.go file like we created myWriter
	// which is just a custom implementation of ResponseWriter
	err = Template(&ww, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	// negative case: when we are using wrong template
	err = Template(&ww, r, "non-existent.page.html", &models.TemplateData{})
	if err == nil {
		t.Error("Expected: Unable to render the template since we are using wrong template name")
	}

	app.UseCache = false
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

func TestCreateTemplateCache(t *testing.T) {

	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
