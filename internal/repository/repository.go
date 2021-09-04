package repository

import "github.com/prayagsingh/bookings/internal/models"

type DatabaseRepo interface {

	// Implemented in postgres.go file
	AllUsers() bool

	InsertReservation(res models.Reservation) (int,error)
	InserRoomRestriction(res models.RoomRestriction) error
}