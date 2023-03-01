package view

import (
	"encoding/json"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
)

// Room a space that can be occupied or where something can be done
// this is a view model for device.Room which abstracts the internal
// implementation details of direction
type Room struct {
	Direction   string `json:"direction"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Room converts the view.Room to device.Room
func (room Room) Room() (gateway.Room, error) {
	direction, err := gateway.NewDirection(room.Direction)
	if err != nil {
		return gateway.Room{}, err
	}

	return gateway.Room{
		Direction: direction,
		PhysicalEntity: gateway.PhysicalEntity{
			Name:        room.Name,
			Description: room.Description,
		},
	}, nil
}

// Convert uses floor and []byte representing view.Room as device.Room
func Convert(floor gateway.Floor, data []byte) (gateway.Room, error) {
	room := Room{}
	err := json.Unmarshal(data, &room)
	if err != nil {
		return gateway.Room{}, err
	}

	r, err := room.Room()
	if err != nil {
		return gateway.Room{}, err
	}

	data, err = json.Marshal(r)
	if err != nil {
		return gateway.Room{}, err
	}
	return gateway.NewRoom(floor, data)
}

// NewRooms convert device.Rooms into []Room
func NewRooms(rooms gateway.Rooms) []Room {
	result := make([]Room, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, NewRoom(room))
	}
	return result
}

// NewRoom converts the view.Room to device.Room
func NewRoom(room gateway.Room) Room {
	return Room{
		Direction:   room.Direction.Direction(),
		Name:        room.Name,
		Description: room.Description,
	}
}
