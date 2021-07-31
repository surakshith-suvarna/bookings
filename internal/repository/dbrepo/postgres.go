package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/surakshith-suvarna/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

//InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	stmt := `Insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

//InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		time.Now(),
		time.Now(),
		r.RestrictionId,
	)

	if err != nil {
		return err
	}
	return nil
}

//SearchAvailabilityByDatesByRoomID return true if availability exists for roomId and false if no availability exists
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `select 
				count(id)
				from room_restrictions
				where
					room_id = $1
					and $2 < end_date and $3 > start_date`

	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms, if any for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select r.id, r.room_name
				from rooms r
				where r.id not in
				(select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date )`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
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

//GetRoomById gets a room by Id
func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}

//GetUserById gets User by ID
func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, created_at, updated_at, access_level
			from users where id=$1`

	var u models.User
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.AccessLevel,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

//UpdateUser updates user record
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1,
			last_name=$2,
			email=$3,
			updated_at=$4,
			access_level=$5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		time.Now(),
		u.AccessLevel)

	if err != nil {
		return err
	}
	return nil
}

//Authenticate checks login credentials
func (m *postgresDBRepo) Authenticate(email, checkpassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashPassword string
	row := m.DB.QueryRowContext(ctx, "select id, password from users where email=$1", email)
	err := row.Scan(&id, &hashPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(checkpassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("invalid password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashPassword, nil
}
