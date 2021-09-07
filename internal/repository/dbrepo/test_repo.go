// This file contains a function which are available to DatabaseRepo interface.
package dbrepo

import (
	"time"

	"github.com/prayagsingh/bookings/internal/models"
)

func (m *testPostgresDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the db
func (m *testPostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

// InserRoomRestriction inserts a room restriction into DB
func (m *testPostgresDBRepo) InserRoomRestriction(res models.RoomRestriction) error {

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exist for roomID else false
func (m *testPostgresDBRepo) SearchAvailabilityByDatesByRoomID(start_date, end_date time.Time, roomID int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of rooms for a given date range
func (m *testPostgresDBRepo) SearchAvailabilityForAllRooms(start_date, end_date time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID get a room by roomID
func (m *testPostgresDBRepo) GetRoomByID(roomID int) (models.Room, error) {

	var room models.Room

	return room, nil
}
