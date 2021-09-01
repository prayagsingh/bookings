package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prayagsingh/bookings/internal/config"
	"github.com/prayagsingh/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(t *testing.M) {

	// storing info to Session
	gob.Register(models.Reservation{})

	// set it to true when in production
	testApp.InProduction = false

	// initialzing logger and printing logs to terminal
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	// make infoLog app wide variable
	testApp.InfoLog = infoLog

	// initializng error logger
	// Lshotfile will give the info about the error
	errorLog := log.New(os.Stdout, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	// make errorLog app wide variable
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// set to true in production
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(t.Run())
}

// creating a http writer to server TestRenderTemplate
type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)

	return length, nil
}
