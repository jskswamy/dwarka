package gateway

import "strings"

// Direction point to or from which a person or thing moves or faces
type Direction int

// Direction returns the string representation of direction
func (direction Direction) Direction() string {
	switch direction {
	case DirectionEast:
		return "east"
	case DirectionSouth:
		return "south"
	case DirectionWest:
		return "west"
	default:
		return "north"
	}
}

// NewDirection converts string direction as Direction
func NewDirection(direction string) Direction {
	switch strings.ToLower(direction) {
	case "east":
		return DirectionEast
	case "south":
		return DirectionSouth
	case "west":
		return DirectionWest
	default:
		return DirectionNorth
	}
}
