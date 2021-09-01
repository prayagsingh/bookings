package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/prayagsingh/bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError
func ClientError(w http.ResponseWriter, status int) {

	app.InfoLog.Println("Client error with status of ", status)
	// showing client the error with status code if something went wrong
	http.Error(w, http.StatusText(status), status)
}

// ServerError
func ServerError(w http.ResponseWriter, err error) {
	
	// getting the stack trace if something went wrong on the server
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
