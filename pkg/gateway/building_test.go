package gateway

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilding_Validate(t *testing.T) {
	t.Run("should not return any validation error", func(t *testing.T) {
		building := Building{
			Lat: 1.2,
			Lan: 1.4,
			PhysicalEntity: PhysicalEntity{
				Name:        "building one",
				Description: "description",
			},
		}

		err := building.Validate()

		assert.NoError(t, err)
	})

	t.Run("should return all the required field validation error", func(t *testing.T) {
		building := Building{
			PhysicalEntity: PhysicalEntity{
				Description: "description",
			},
		}

		err := building.Validate()

		if assert.Error(t, err) {
			assert.Equal(t, "lan: cannot be blank; lat: cannot be blank; name: cannot be blank.", err.Error())
		}
	})

	t.Run("should error if name provided is not of valid length", func(t *testing.T) {
		building := Building{
			Lat: 1.2,
			Lan: 1.4,
			PhysicalEntity: PhysicalEntity{
				Name:        "name",
				Description: "description",
			},
		}

		err := building.Validate()

		if assert.Error(t, err) {
			assert.Equal(t, "name: the length must be between 5 and 50.", err.Error())
		}
	})
}

func TestNewBuilding(t *testing.T) {
	t.Run("it should create new building", func(t *testing.T) {
		building := Building{
			Lat: 1.2,
			Lan: 1.4,
			PhysicalEntity: PhysicalEntity{
				Name:        "building-one",
				Description: "description",
			},
		}
		data, _ := json.Marshal(building)

		actual, err := NewBuilding(data)

		assert.NoError(t, err)
		assert.Equal(t, building, actual)
	})

	t.Run("return validation error if any", func(t *testing.T) {
		building := Building{
			Lat: 1.2,
			Lan: 1.4,
			PhysicalEntity: PhysicalEntity{
				Name:        "name",
				Description: "description",
			},
		}
		data, _ := json.Marshal(building)

		actual, err := NewBuilding(data)

		if assert.Error(t, err) {
			assert.Equal(t, "name: the length must be between 5 and 50.", err.Error())
		}
		assert.Equal(t, building, actual)
	})
}

func TestNewBuildings(t *testing.T) {
	t.Run("it should create new buildings", func(t *testing.T) {
		buildings := Buildings{"building-one": Building{
			Lat: 1.2,
			Lan: 1.4,
			PhysicalEntity: PhysicalEntity{
				Name:        "building-one",
				Description: "description",
			},
		}}
		data, _ := json.Marshal(buildings)

		actual, err := NewBuildings(data)

		assert.NoError(t, err)
		assert.Equal(t, buildings, actual)
	})
}
