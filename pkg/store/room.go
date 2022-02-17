package store

import (
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"path"
)

const (
	roomsBasePath = "rooms"
)

// Rooms returns all the Rooms from store
func (ps PersistentStore) Rooms(floor gateway.Floor) (gateway.Rooms, error) {
	value, err := ps.get(ps.roomsRootPath(floor), gateway.Rooms{})
	if err != nil {
		return nil, err
	}
	return gateway.NewRooms(floor, value)
}

// UpsertRooms creates or updates Rooms in store
func (ps PersistentStore) UpsertRooms(floor gateway.Floor, rooms gateway.Rooms) error {
	return ps.putJSON(ps.roomsRootPath(floor), rooms)
}

// UpsertRoom creates or updates Room in store
func (ps PersistentStore) UpsertRoom(room gateway.Room) error {
	rooms, err := ps.Rooms(room.Floor)
	if err != nil {
		return err
	}
	rooms[room.ID()] = room
	return ps.putJSON(ps.roomsRootPath(room.Floor), rooms)
}

// DeleteRoom deletes the room and nested path from store
func (ps PersistentStore) DeleteRoom(room gateway.Room) error {
	rooms, err := ps.Rooms(room.Floor)
	if err != nil {
		return err
	}

	delete(rooms, room.ID())
	err = ps.putJSON(ps.roomsRootPath(room.Floor), rooms)
	if err != nil {
		return err
	}

	return ps.safeDelete(ps.roomRootPath(room))
}

func (ps PersistentStore) roomsRootPath(floor gateway.Floor) string {
	return path.Join(ps.floorRootPath(floor), roomsBasePath)
}

func (ps PersistentStore) roomRootPath(room gateway.Room) string {
	return path.Join(ps.floorRootPath(room.Floor), room.ID())
}
