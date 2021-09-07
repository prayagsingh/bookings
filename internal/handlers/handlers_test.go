package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/prayagsingh/bookings/internal/models"
)

// for sending data for POST request
type postData struct {
	key   string // key is the name of the form element or form input
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	// {"home", "/", "GET", []postData{}, http.StatusOK},
	// {"about", "/about", "GET", []postData{}, http.StatusOK},
	// {"villas", "/villas", "GET", []postData{}, http.StatusOK},
	// {"suites", "/suites", "GET", []postData{}, http.StatusOK},
	// {"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	// {"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// {"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},

	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{key:"start", value: "2021-10-01"},
	// 	{key:"end", value: "2021-10-02"},
	// }, http.StatusOK},

	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key:"start", value: "2021-10-01"},
	// 	{key:"end", value: "2021-10-02"},
	// }, http.StatusOK},

	// {"make-reservation-post", "/make-reservation", "POST", []postData{
	// 	{key:"first_name", value: "Leo"},
	// 	{key:"last_name", value: "Messi"},
	// 	{key:"email", value: "leomessi@psg.com"},
	// 	{key:"phone", value: "111-111-1111"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()

	// creating a http test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		// what kind of request we are sending
		if e.method == "GET" {
			// we need to make a request as a client here
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			// if no error then check for the status codes
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else {
			// POST requests
			// create a variable that is in the format that our testserver is expecting
			// creating an empty variable that is in the format which is required by the method that we gonna call on our test-server
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Error(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {

	// we need to have models.Reservation to put into the session
	// here using Villas for test hence roomID is 1
	reservation := models.Reservation{
		RoomID: 1, // for villas
		Room: models.Room{
			ID:       1,
			RoomName: "Villas",
		},
	}

	request, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(request)
	request = request.WithContext(ctx)

	// NewRecorder simulates what we get from the request response lifecycle
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservations)
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler retured wrong status code. expected %d and got %d", http.StatusOK, rr.Code)
	}
}

// helper func for putting the reservation var as a session var into the session of the request
// and it is possible using context hence creating a getCtx helper func
func getCtx(r *http.Request) context.Context {

	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
