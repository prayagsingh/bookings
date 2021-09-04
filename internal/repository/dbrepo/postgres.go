// This file contains a function which are available to DatabaseRepo interface.
package dbrepo

import (
	"context"
	"time"

	"github.com/prayagsingh/bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the db
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	// creating context to make sure that the txn should not open for more than set time like adding a default timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// need reservation-id for room-restriction table
	var newID int

	stmt := `insert into reservations (first_name , last_name, email, phone, start_date,
	        end_date, room_id, created_at, updated_at)
			values($1, $2,$3, $4, $5, $6, $7, $8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

// InserRoomRestriction inserts a room restriction into DB
func (m *postgresDBRepo) InserRoomRestriction(res models.RoomRestriction) error {

	// creating context to make sure that the txn should not open for more than set time like adding a default timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id,
		    created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
