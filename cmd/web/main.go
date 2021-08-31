package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prayagsingh/bookings/internal/config"
	"github.com/prayagsingh/bookings/internal/handlers"
	"github.com/prayagsingh/bookings/internal/models"
	"github.com/prayagsingh/bookings/internal/render"
)

const portNumber = ":8080"

// making AppConfig available to all the files under package main
var app config.AppConfig

// making session available to all the files under package main
var session *scs.SessionManager

func main() {

	// storing info to Session
	gob.Register(models.Reservation{})

	// set it to true when in production
	app.InProduction = false

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

	fmt.Printf(fmt.Sprintf("Starting application on port %s\n", portNumber))

	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	// dirPath, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("\ncurrent directory is: ", dirPath)

	// execPath, err := os.Executable()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("\n executable path is: ", execPath)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
