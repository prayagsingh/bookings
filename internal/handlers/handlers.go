package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prayagsingh/bookings/internal/config"
	"github.com/prayagsingh/bookings/internal/forms"
	"github.com/prayagsingh/bookings/internal/models"
	"github.com/prayagsingh/bookings/internal/render"
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
	//fmt.Printf("\nrequest url path is: %s\n", r.URL.Path)
	render.RenderTemplate(rw, r, "home.page.html", &models.TemplateData{})

}

// About is the handler for the about page
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World !!!"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	remotePort := m.App.Session.GetString(r.Context(), "remote_port")
	stringMap["remote_ip"] = remoteIp
	stringMap["remote_port"] = remotePort

	render.RenderTemplate(rw, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservations renders a make a reservation page and displays form
func (m *Repository) Reservations(rw http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(rw, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservations handles the posting of a reservation form
func (m *Repository) PostReservations(rw http.ResponseWriter, r *http.Request) {

	// `name` attribute must be present in `input` element to fetch the values from Form
	err := r.ParseForm()
	if err != nil {
		log.Println("Error in parsing form data and error is ", err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	// creating a form object to check our data
	form := forms.New(r.PostForm)

	// using below in make-reservation page for showing warnings
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)

	if !form.Valid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(rw, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

}

// Villas renders the room page
func (m *Repository) Villas(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, r, "villas.page.html", &models.TemplateData{})
}

// Suites renders the room page
func (m *Repository) Suites(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, r, "suites.page.html", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {
	// fetching values from the form using ID's mentioned in the form
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	rw.Write([]byte(fmt.Sprintf("Start date is %s and End date is %s", start, end)))
}

// AvailabiltyJSON is using it to build JSON response. Scope is limited
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response.
func (m *Repository) AvailabilityJSON(rw http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Println("Error in marshelling and error is: ", err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}

// Contact renders the search contact page
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, r, "contact.page.html", &models.TemplateData{})
}
