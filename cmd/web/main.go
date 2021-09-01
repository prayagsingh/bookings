package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prayagsingh/bookings/internal/config"
	"github.com/prayagsingh/bookings/internal/handlers"
	"github.com/prayagsingh/bookings/internal/helpers"
	"github.com/prayagsingh/bookings/internal/models"
	"github.com/prayagsingh/bookings/internal/render"
)

const portNumber = ":8080"

// making AppConfig available to all the files under package main
var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

// making session available to all the files under package main
var session *scs.SessionManager

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(fmt.Sprintf("Starting application on port %s\n", portNumber))

	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// storing info to Session
	gob.Register(models.Reservation{})

	// set it to true when in production
	app.InProduction = false

	// initialzing logger and printing logs to terminal
	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	// make infoLog app wide variable
	app.InfoLog = infoLog

	// initializng error logger
	// Lshotfile will give the info about the error
	errorLog = log.New(os.Stdout, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	// make errorLog app wide variable
	app.ErrorLog = errorLog

	// Intializing a SessionManager
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// set to true in production
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("can't create template cache")
		return err
	}

	app.TemplateCache = tc
	// set to false if you want to load template from disk in development mode
	// if it is set to true then any changes made to templates won't reflect dynamically because
	// it is reading from template cache instead of disk. it is faster in loading when comparing it
	// to reading it from disk
	app.UseCache = false

	// This allow Handler functions to have access to appConfig via repository
	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplate(&app)
	helpers.NewHelpers(&app)

	return nil
}
