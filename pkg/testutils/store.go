package testutils

import (
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"time"
)

// NewBuilding return new building from name
func NewBuilding(name string) gateway.Building {
	return gateway.Building{
		PhysicalEntity: gateway.PhysicalEntity{Name: name},
	}
}

// NewBuildings creates and returns a building from name and
// buildings after associating it
func NewBuildings(name string) (gateway.Buildings, gateway.Building) {
	building := NewBuilding(name)
	return gateway.Buildings{name: building}, building
}

// NewFloor return new floor from name
func NewFloor(name string) gateway.Floor {
	return gateway.Floor{
		Building:       NewBuilding("building-one"),
		PhysicalEntity: gateway.PhysicalEntity{Name: name},
	}
}

// NewFloors creates and returns a floor from name and
// floors after associating it
func NewFloors(name string) (gateway.Floors, gateway.Floor) {
	floor := NewFloor(name)
	return gateway.Floors{name: floor}, floor
}

// AssociateFloorToBuilding associate the floor to the building
func AssociateFloorToBuilding(building gateway.Building, floor gateway.Floor) gateway.Floor {
	floor.Building = building
	return floor
}

// Uptime return current time wrapped as status
func Uptime() gateway.Status {
	now := time.Now().Local().Format(time.RFC822)
	return gateway.Status{"startTime": now}
}
