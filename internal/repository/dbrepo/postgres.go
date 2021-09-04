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

// SearchAvailabilityByDatesByRoomID returns true if availability exist for roomID else false
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start_date, end_date time.Time, roomID int) (bool, error) {

	// creating context to make sure that the txn should not open for more than set time like adding a default timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select
			count(id)
		from
			room_restrictions
		where 
			room_id = $1
		and	
			$2 < end_date 
		and 
			$3 > start_date;
	`
	var numRows int
	err := m.DB.QueryRowContext(ctx, query, roomID, start_date, end_date).Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of rooms for a given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start_date, end_date time.Time) ([]models.Room, error) {

	// creating context to make sure that the txn should not open for more than set time like adding a default timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	var rooms []models.Room
	query := `select
	 	r.id, r.room_name
	from
		rooms r
	where
		r.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);
	`
	// Here we are querying multiple rows hence using QueryContext instead of QueryRowContext
	rows, err := m.DB.QueryContext(ctx, query, start_date, end_date)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(
			&room.ID, 
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}