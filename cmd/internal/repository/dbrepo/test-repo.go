package dbrepo

import (
	"errors"
	"time"

	"github.com/maslow123/bookings/cmd/internal/models"
)

func (m *testDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InserRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error { // if the room id is 2, then fail; otherwise, pass
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms return  a slice of available rooms, if any. For given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID return room data
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("Some error")
	}

	return room, nil
}

// GetUserByID return user data
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User

	return user, nil
}

// UpdateUser updates a user in the database
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate autheticates a user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil

}

// AllNewReservations returns a slice of all reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation
	return reservations, nil

}

// GetReservationByID returns one reservation by ID
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation

	return res, nil
}

// UpdateReservation updates a reservation in the database
func (m *testDBRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

// DeleteReservation deletes one reservations by id
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}
