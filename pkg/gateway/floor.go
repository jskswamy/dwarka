package gateway

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

// Floor a horizontal plane or line with respect to the distance above or below a given point
type Floor struct {
	Building Entity `json:"-"`
	Level    int    `json:"level"`
	PhysicalEntity
}

// Validate validate whether floor has all the necessary fields
func (floor Floor) Validate() error {
	return validation.ValidateStruct(&floor,
		validation.Field(&floor.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&floor.Level, validation.Required),
	)
}

// NewFloor returns a Floor from []byte
func NewFloor(building Building, data []byte) (Floor, error) {
	floor := Floor{Building: building}
	err := json.Unmarshal(data, &floor)
	if err != nil {
		return Floor{}, fmt.Errorf("unable to parse floor, %w", err)
	}

	err = floor.Validate()
	if err != nil {
		return floor, err
	}

	return floor, nil
}

// Floors represents map string, Floor
type Floors map[string]Floor

// NewFloors returns list of Floors from []byte
func NewFloors(building Entity, data []byte) (Floors, error) {
	floors := Floors{}
	err := json.Unmarshal(data, &floors)
	if err != nil {
		return nil, fmt.Errorf("unable to parse floors, %w", err)
	}

	result := Floors{}

	for _, floor := range floors {
		floor.Building = building
		result[floor.ID()] = floor
	}
	return result, nil
}
