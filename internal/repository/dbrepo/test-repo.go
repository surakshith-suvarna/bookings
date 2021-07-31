package dbrepo

import (
	"errors"
	"time"

	"github.com/surakshith-suvarna/bookings/internal/models"
)

//InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//If roomID is 2 then fail otherwise pass
	if res.RoomId == 2 {
		return 0, errors.New("insertion failed")
	}

	return 0, nil
}

//InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomId == 1000 {
		return errors.New("insert restrictions failed")
	}
	return nil
}

//SearchAvailabilityByDatesByRoomID return true if availability exists for roomId and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {

	return false, nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms, if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

//GetRoomById gets a room by Id
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("roomID not found")
	}

	return room, nil
}

func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, checkpassword string) (int, string, error) {
	return 1, "", nil
}
