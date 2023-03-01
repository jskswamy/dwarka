package gateway

import (
	"errors"
	"fmt"
	"strings"
)

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
func NewDirection(direction string) (Direction, error) {
	switch strings.ToLower(direction) {
	case "east":
		return DirectionEast, nil
	case "south":
		return DirectionSouth, nil
	case "west":
		return DirectionWest, nil
	case "north":
		return DirectionNorth, nil
	default:
		return -1, errors.New(fmt.Sprintf("direction %s not supported", direction))
	}
}
