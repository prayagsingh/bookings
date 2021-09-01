package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/prayagsingh/bookings/internal/config"
	"github.com/prayagsingh/bookings/internal/models"
	"github.com/prayagsingh/bookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {

	// storing info to Session
	gob.Register(models.Reservation{})

	// set it to true when in production
	app.InProduction = false

	// initialzing logger and printing logs to terminal
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	// make infoLog app wide variable
	app.InfoLog = infoLog

	// initializng error logger
	// Lshotfile will give the info about the error
	errorLog := log.New(os.Stdout, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	// make errorLog app wide variable
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// set to true in production
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("can't create template cache")

	}

	app.TemplateCache = tc
	// set to false if you want to load template from disk in development mode
	// if it is set to true then any changes made to templates won't reflect dynamically because
	// it is reading from template cache instead of disk. it is faster in loading when comparing it
	// to reading it from disk
	app.UseCache = true

	// This allow Handler functions to have access to appConfig via repository
	repo := NewRepo(&app)
	NewHandler(repo)

	render.NewTemplate(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	// this will return BAD request if any request don't have a valid csrf token
	// not required in testcases
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/villas", Repo.Villas)
	mux.Get("/suites", Repo.Suites)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservations)
	mux.Post("/make-reservation", Repo.PostReservations)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// routes for static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf add CSRF protection to all the POST requests
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		// set it to true when in production
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	// store the template found during parsing
	myCache := map[string]*template.Template{}

	// find all the pages in template dir which ends with page.html
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

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
	}

	return myCache, nil
}
