package gateway

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

// Building  a structure with a roof and walls, such as a house or factory
type Building struct {
	Lat float64 `json:"lat"`
	Lan float64 `json:"lan"`
	PhysicalEntity
}

// Validate validate whether building has all the necessary fields
func (building Building) Validate() error {
	return validation.ValidateStruct(&building,
		validation.Field(&building.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&building.Lat, validation.Required),
		validation.Field(&building.Lan, validation.Required),
	)
}

// NewBuilding returns a Building from []byte
func NewBuilding(data []byte) (Building, error) {
	building := Building{}
	err := json.Unmarshal(data, &building)
	if err != nil {
		return Building{}, fmt.Errorf("unable to parse building, %w", err)
	}

	err = building.Validate()
	if err != nil {
		return building, err
	}

	return building, nil
}

// Buildings represents map string, Building
type Buildings map[string]Building

// NewBuildings returns list of Buildings from []byte
func NewBuildings(data []byte) (Buildings, error) {
	buildings := Buildings{}
	err := json.Unmarshal(data, &buildings)
	if err != nil {
		return nil, fmt.Errorf("unable to parse buildings, %w", err)
	}

	return buildings, nil
}
