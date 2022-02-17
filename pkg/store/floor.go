package store

import (
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"path"
)

const (
	floorsBasePath = "floors"
)

// Floors returns all the Floors from store
func (ps PersistentStore) Floors(building gateway.Entity) (gateway.Floors, error) {
	value, err := ps.get(ps.floorsRootPath(building), gateway.Floors{})
	if err != nil {
		return nil, err
	}
	return gateway.NewFloors(building, value)
}

// UpsertFloors creates or updates Floors in store
func (ps PersistentStore) UpsertFloors(building gateway.Entity, floors gateway.Floors) error {
	return ps.putJSON(ps.floorsRootPath(building), floors)
}

// UpsertFloor creates or updates Floor in store
func (ps PersistentStore) UpsertFloor(floor gateway.Floor) error {
	floors, err := ps.Floors(floor.Building)
	if err != nil {
		return err
	}
	floors[floor.ID()] = floor
	return ps.putJSON(ps.floorsRootPath(floor.Building), floors)
}

// DeleteFloor deletes the floor and nested path from store
func (ps PersistentStore) DeleteFloor(floor gateway.Floor) error {
	floors, err := ps.Floors(floor.Building)
	if err != nil {
		return err
	}

	delete(floors, floor.ID())
	err = ps.putJSON(ps.floorsRootPath(floor.Building), floors)
	if err != nil {
		return err
	}

	return ps.safeDelete(ps.floorRootPath(floor))
}

func (ps PersistentStore) floorsRootPath(building gateway.Entity) string {
	return path.Join(ps.buildingRootPath(building), floorsBasePath)
}

func (ps PersistentStore) floorRootPath(floor gateway.Floor) string {
	return path.Join(ps.buildingRootPath(floor.Building), floor.ID())
}
