package repository

import (
	"time"

	"github.com/surakshith-suvarna/bookings/internal/models"
)

type DatabaseRepo interface {
	//Instead of putting all reservation model items one by one we are declaring as models
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)

	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, checkpassword string) (int, string, error)
}
