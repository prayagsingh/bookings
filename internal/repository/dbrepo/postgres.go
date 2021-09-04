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
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {

	// creating context to make sure that the txn should not open for more than set time like adding a default timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservations (first_name , last_name, email, phone, start_date,
	        end_date, room_id, created_at, updated_at)
			values($1, $2,$3, $4, $5, $6, $7, $8,$9)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
