package repository

import (
	"time"

	"github.com/prayagsingh/bookings/internal/models"
)

type DatabaseRepo interface {

	// Implemented in postgres.go file
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InserRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start_date, end_date time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start_date, end_date string) ([]models.Room, error)
}
