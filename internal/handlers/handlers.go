// casting start_date, end_date of type string to time.Time
// ref: 01/02 03:04:05PM '06 -0700
// article ref: https://www.pauladamsmith.com/blog/2011/05/go_time.html
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/prayagsingh/bookings/internal/config"
	"github.com/prayagsingh/bookings/internal/driver"
	"github.com/prayagsingh/bookings/internal/forms"
	"github.com/prayagsingh/bookings/internal/helpers"
	"github.com/prayagsingh/bookings/internal/models"
	"github.com/prayagsingh/bookings/internal/render"
	"github.com/prayagsingh/bookings/internal/repository"
	"github.com/prayagsingh/bookings/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	// making sure that the db is available to handlers
	DB repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {

	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository for testcases
func NewTestRepo(a *config.AppConfig) *Repository {

	return &Repository{
		App: a,
		DB:  dbrepo.NewTestPostgresRepo(a),
	}
}

// NewHandler sets the repository for the handler
func NewHandler(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "home.page.html", &models.TemplateData{})

}

// About is the handler for the about page
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "about.page.html", &models.TemplateData{})
}

// Reservations renders a make a reservation page and displays form
func (m *Repository) Reservations(rw http.ResponseWriter, r *http.Request) {

	// if we try to access /make-reservation directly then it will result in session error
	// because earlier we were creating an empty reservation
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(rw, errors.New("can't get reservation from session"))
		return
	}

	// get the room by room-id to display it on the page
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	// storing room name to reservation
	res.Room.RoomName = room.RoomName

	// putting room name to session
	m.App.Session.Put(r.Context(), "reservation", res)

	// converting time.Time formart to string to serve to html page
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	// passing it as a stringmap in templatedata so that html page can have access
	dates := make(map[string]string)
	dates["start_date"] = sd
	dates["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(rw, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		StringMap: dates,
		Data:      data,
	})
}

// PostReservations handles the posting of a reservation form
func (m *Repository) PostReservations(rw http.ResponseWriter, r *http.Request) {

	// pulling out the room name and other details from the session
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(rw, errors.New("can't get from session"))
		return
	}

	// `name` attribute must be present in `input` element to fetch the values from Form
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	// creating a form object to check our data
	form := forms.New(r.PostForm)

	// using below in make-reservation page for showing warnings
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(rw, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// putting reservation data to DB
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	restriction := models.RoomRestriction{
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
	}

	err = m.DB.InserRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	// showing the reservation summary using session.  to do this we have to pass the reservation
	// object to session and when we get to reservation-sumary page then we will pull out the object
	// from Session and finally sent it to the template and display the information
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// To avoid the people accidently submit the form twice, any time we recieve the POST request
	// we should directs the user to another page with a HTTP redirect 303
	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)

}

// Villas renders the room page
func (m *Repository) Villas(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "villas.page.html", &models.TemplateData{})
}

// Suites renders the room page
func (m *Repository) Suites(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "suites.page.html", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {
	// fetching values from the form using ID's mentioned in the form
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	if len(rooms) == 0 {
		// no room available
		m.App.Session.Put(r.Context(), "error", "No rooms available !!!")
		// redirecting with 303 status code
		http.Redirect(rw, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(rw, r, "choose-rooms.page.html", &models.TemplateData{
		Data: data,
	})
}

// AvailabiltyJSON is using it to build JSON response. Scope is limited
type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and sends JSON response.
func (m *Repository) AvailabilityJSON(rw http.ResponseWriter, r *http.Request) {

	// fetching start_date and end_date from the form
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	// convert dates to time.Time format from string
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(rw, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	// making json resp dynamic based on the response we get from the DB .i.e whether room is available for not
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}

// Contact renders the search contact page
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "contact.page.html", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {

	// doing type assert(added models.Reservation) to identify what type of session it is.
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("get error from session")
		// if a user directly went to /reservation-summary page directly then it will show empty page
		// because of lack of session hence we have to show them something if they directly went to
		// reservation-summary page
		m.App.Session.Put(r.Context(), "error", "can't get reservation from the session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)

		return
	}

	// removing the reservation from session
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(rw, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {

	// we need to get the ID from chi router
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	//m.App.InfoLog.Println("RoomID is: ", roomID)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	// now need to get the reservarion variable stored in session then update the room-id and
	// put the reservation back to session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(rw, err)
		return
	}

	res.RoomID = roomID
	// putting back the reservation to session after updating the room-id
	m.App.Session.Put(r.Context(), "reservation", res)

	// redirecting to make-reservation page
	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL parameter, build a sessional variable and take user to reservation screen
func (m *Repository) BookRoom(rw http.ResponseWriter, r *http.Request) {

	// fetch the id, sd, ed from url
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(rw, err)
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	// storing above details in session
	var reservation models.Reservation

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	// get the room by room-id to display it on the page
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(rw, err)
	}

	reservation.RoomID = roomID
	reservation.Room.RoomName = room.RoomName
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)
}
