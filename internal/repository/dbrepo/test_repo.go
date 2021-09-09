// This file contains a function which are available to DatabaseRepo interface.
package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/prayagsingh/bookings/internal/models"
)

func (m *testPostgresDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the db
func (m *testPostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2 then fail otherwise pass
	if res.RoomID == 2 {
		return 0, errors.New("invalid room ID")
	}
	return 1, nil
}

// InserRoomRestriction inserts a room restriction into DB
func (m *testPostgresDBRepo) InserRoomRestriction(res models.RoomRestriction) error {
	// for room id 1000, make it faile
	if res.RoomID == 1000 {
		return errors.New("room id is incorrect")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exist for roomID else false
func (m *testPostgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	// set up a test time
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start.String() == testDateToFail.String() {
		return false, errors.New("some error")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}
	// otherwise, we have availability
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of rooms for a given date range
func (m *testPostgresDBRepo) SearchAvailabilityForAllRooms(start_date, end_date time.Time) ([]models.Room, error) {

	var rooms []models.Room

	// if the start date is after 2049-12-31, then return empty slice,
	// indicating no rooms are available;
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start_date == testDateToFail {
		return rooms, errors.New("some error")
	}

	if start_date.After(t) {
		return rooms, nil
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

// GetRoomByID get a room by roomID
func (m *testPostgresDBRepo) GetRoomByID(roomID int) (models.Room, error) {

	var room models.Room

	if roomID > 2 {
		return room, errors.New("room not found. room-id is greator than 2")
	}
	return room, nil
}
