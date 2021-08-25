package handlers

import (
	"log"
	"net"
	"net/http"

	"github.com/prayagsingh/bookings/pkg/config"
	"github.com/prayagsingh/bookings/pkg/models"
	"github.com/prayagsingh/bookings/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {

	return &Repository{
		App: a,
	}
}

// NewHandler sets the repository for the handler
func NewHandler(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	remoteIP, remotePort, err := net.SplitHostPort(remoteIP)
	if err != nil {
		log.Fatalf("Unable to fetch the remoteIP and error is: %s", err)
	}
	// storing the remote_IP to session
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	m.App.Session.Put(r.Context(), "remote_port", remotePort)
	render.RenderTemplate(rw, "home.page.html", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World !!!"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	remotePort := m.App.Session.GetString(r.Context(), "remote_port")
	stringMap["remote_ip"] = remoteIp
	stringMap["remote_port"] = remotePort

	render.RenderTemplate(rw, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
