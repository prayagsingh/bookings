package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/prayagsingh/bookings/internal/config"
	"github.com/prayagsingh/bookings/internal/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// testcases can access the templates from root folder
var pathToTemplates = "./templates"

// NewTemplate sets the config for the new template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all the templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	// if a user directly went to /reservation-summary page directly then it will show empty page
	// because of lack of session hence we have to show them something if they directly went to
	// reservation-summary page
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)

	return td
}

// RenderTemplate for rendering the template using html/template
func RenderTemplate(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var templateCache map[string]*template.Template
	// In Production load template from template Cache
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}
	// ok is used to check if the template exists or not
	// if template found, ok wil return true else false
	t, ok := templateCache[tmpl]

	//fmt.Println("t is ", t.DefinedTemplates())

	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		fmt.Println("Error writing parsed template to buffer", err)
		return nil
	}
	_, err = buf.WriteTo(rw)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return nil
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	// store the template found during parsing
	myCache := map[string]*template.Template{}

	// find all the pages in template dir which ends with page.html
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Printf("Page is currently %s and name is %s \n", page, name)
		// New allocates the new template with a given name
		// Funcs adds the elements of the argument map to the template's function map. It must be
		// 	 called before the template is parsed.
		// ParseFiles parses the named files and associates the resulting templates with t. If an
		//   error occurs, parsing stops and the returned template is nil otherwise it is t
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		// Adding templates to cache
		myCache[name] = ts
		//fmt.Println("Defined templates stored in myCache is ", myCache[name].DefinedTemplates())
		//fmt.Println("template stored in myCache is ", myCache[name].Tree.Name)
	}

	return myCache, nil
}
