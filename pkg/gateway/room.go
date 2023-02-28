package gateway

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

// Room a space that can be occupied or where something can be done
type Room struct {
	Floor     Floor `json:"-"`
	Direction `json:"direction"`
	PhysicalEntity
}

// Validate validates whether room has all the necessary fields
func (room Room) Validate() error {
	return validation.ValidateStruct(&room,
		validation.Field(&room.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&room.Direction, validation.Required),
	)
}

// NewRoom returns a Room from []byte
func NewRoom(floor Floor, data []byte) (Room, error) {
	room := Room{Floor: floor}
	err := json.Unmarshal(data, &room)
	if err != nil {
		return Room{}, fmt.Errorf("unable to parse room, %w", err)
	}

	err = room.Validate()
	if err != nil {
		return room, err
	}

	return room, nil
}

// Rooms represents map string, Room
type Rooms map[string]Room

// NewRooms returns list of Rooms from []byte
func NewRooms(floor Floor, data []byte) (Rooms, error) {
	rooms := Rooms{}
	err := json.Unmarshal(data, &rooms)
	if err != nil {
		return nil, fmt.Errorf("unable to parse rooms, %w", err)
	}

	result := Rooms{}

	for _, room := range rooms {
		room.Floor = floor
		result[room.ID()] = room
	}
	return result, nil
}
