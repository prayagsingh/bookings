package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

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
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"villas", "/villas", "GET", http.StatusOK},
	{"suites", "/suites", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	//{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},

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

			// } else {
			// 	// POST requests
			// 	// create a variable that is in the format that our testserver is expecting
			// 	// creating an empty variable that is in the format which is required by the method that we gonna call on our test-server
			// 	values := url.Values{}

			// 	for _, x := range e.params {
			// 		values.Add(x.key, x.value)
			// 	}

			// 	resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			// 	if err != nil {
			// 		t.Log(err)
			// 		t.Error(err)
			// 	}

			// 	if resp.StatusCode != e.expectedStatusCode {
			// 		t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			// 	}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {

	// Case 1: when reservation is in session
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

	// Case 2: when reservation is not in session(reset everything)
	request, _ = http.NewRequest("GET", "/make-reservation", nil)
	// If we don't do this then won't be able to find the value in session because there is no session
	ctx = getCtx(request)
	request = request.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, request)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler retured wrong status code. expected %d and got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// Case 3: test with non-existent room
	// testcase when reservation is not in session(reset everything)
	request, _ = http.NewRequest("GET", "/make-reservation", nil)
	// If we don't do this then won't be able to find the value in session because there is no session
	ctx = getCtx(request)
	request = request.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, request)
	// setting room-id to 3 in reservation for this testcase
	reservation.RoomID = 3
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, request)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler retured wrong status code. expected %d and got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	// Case 1: when reservation is in session
	// we need to have models.Reservation to put into the session
	// here using Villas for test hence roomID is 1
	layout := "2006-01-02"
	sd, _ := time.Parse(layout, "2021-01-02")
	ed, _ := time.Parse(layout, "2021-01-03")

	reservation := models.Reservation{
		RoomID:    1, // for villas
		StartDate: sd,
		EndDate:   ed,
		Room: models.Room{
			ID:       1,
			RoomName: "Villas",
		},
	}

	postedData := url.Values{}
	postedData.Add("first_name", "Prayag")
	postedData.Add("last_name", "Singh")
	postedData.Add("email", "ps@myemail.com")
	postedData.Add("phone", "1111111111")
	postedData.Add("room_id", "1")

	request, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(request)
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)
	// NewRecorder simulates what we get from the request response lifecycle
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservations)
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for missing form body
	request, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservations)

	handler.ServeHTTP(rr, request)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong status code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test Form isInvalid
	postedData = url.Values{}
	postedData.Add("first_name", "aa")
	postedData.Add("last_name", "bbb")
	postedData.Add("room_id", "1")

	request, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservations)

	handler.ServeHTTP(rr, request)
	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returned wrong status code for form invalid: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// Test when reservation is not putted into session
	request, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservations)

	handler.ServeHTTP(rr, request)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong status code for case when reservation not putted into session: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Case when unable to insert reservation into DB
	postedData = url.Values{}
	postedData.Add("first_name", "Prayag")
	postedData.Add("last_name", "Singh")
	postedData.Add("email", "ps@email.com")
	postedData.Add("phone", "1111111111")
	postedData.Add("room_id", "1")

	request, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	// sending room id = 2 to make it fail
	reservation.RoomID = 2

	session.Put(ctx, "reservation", reservation)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservations)

	handler.ServeHTTP(rr, request)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for unable to insert reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test when unable to insert room restrictions into DB
	postedData = url.Values{}
	postedData.Add("first_name", "Prayag")
	postedData.Add("last_name", "Singh")
	postedData.Add("email", "ps@email.com")
	postedData.Add("phone", "1111111111")
	postedData.Add("room_id", "1")

	request, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	// sending roomID 1000 to make it test the failure condition
	reservation.RoomID = 1000

	session.Put(ctx, "reservation", reservation)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservations)

	handler.ServeHTTP(rr, request)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	/*****************************************
	// first case -- rooms are not available
	*****************************************/
	// create our request body
	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-02")
	postedData.Add("room_id", "1")

	// create our request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	// this time we want to parse JSON and get the expected response
	var j jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date > 2049-12-31, we expect no availability
	if j.OK {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}

	/*****************************************
	// second case -- rooms not available
	*****************************************/
	// create our request body
	// create our request body
	postedData = url.Values{}
	postedData.Add("start", "2040-01-01")
	postedData.Add("end", "2040-01-02")
	postedData.Add("room_id", "1")

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if !j.OK {
		t.Error("Got no availability when some was expected in AvailabilityJSON")
	}

	/*****************************************
	// third case -- no request body
	*****************************************/
	// create our request
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if j.OK || j.Message != "Internal server error" {
		t.Error("Got availability when request body was empty")
	}

	/*****************************************
	// fourth case -- database error
	*****************************************/
	// create our request body
	// create our request body
	postedData = url.Values{}
	postedData.Add("start", "2060-01-01")
	postedData.Add("end", "2060-01-02")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date > 2049-12-31, we expect availability
	if j.OK || j.Message == "Error querying database" {
		t.Error("Got availability when simulating database error")
	}
}
func TestRepository_PostAvailability(t *testing.T) {
	/*****************************************
	// first case -- rooms are not available
	*****************************************/
	// create our request body
	// create our request body
	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-02")

	// create our request
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// second case -- rooms are available
	*****************************************/
	// this time, we specify a start date before 2040-01-01, which will give us
	// a non-empty slice, indicating that rooms are available
	// create our request body
	postedData = url.Values{}
	postedData.Add("start", "2040-01-01")
	postedData.Add("end", "2040-01-02")

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusOK
	if rr.Code != http.StatusOK {
		t.Errorf("Post availability when rooms are available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	/*****************************************
	// third case -- empty post body
	*****************************************/
	// create our request with a nil body, so parsing form fails
	req, _ = http.NewRequest("POST", "/search-availability", nil)

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	/*****************************************
	// fourth case -- start date in wrong format
	*****************************************/
	// this time, we specify a start date in the wrong format
	// create our request body
	postedData = url.Values{}
	postedData.Add("start", "invalid")
	postedData.Add("end", "2040-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	/*****************************************
	// fifth case -- end date in wrong format
	*****************************************/
	// this time, we specify a start date in the wrong format
	// create our request body
	postedData = url.Values{}
	postedData.Add("start", "2040-01-01")
	postedData.Add("end", "invalid")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	/*****************************************
	// sixth case -- database query fails
	*****************************************/
	// this time, we specify a start date of 2060-01-01, which will cause
	// our testdb repo to return an error
	postedData = url.Values{}
	postedData.Add("start", "2060-01-01")
	postedData.Add("end", "2060-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	/*****************************************
	// first case -- reservation in session
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	/*****************************************
	// second case -- reservation not in session
	*****************************************/
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	/*****************************************
	// first case -- reservation in session
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/choose-room/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	///*****************************************
	//// second case -- reservation not in session
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	///*****************************************
	//// third case -- missing url parameter, or malformed parameter
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/fish", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/fish"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	/*****************************************
	// first case -- database works
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// second case -- database failed
	*****************************************/
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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
