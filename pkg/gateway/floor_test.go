package gateway_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/testutils"
	"testing"
)

func TestNewFloor(t *testing.T) {
	t.Run("should return floor associated to a building", func(t *testing.T) {
		building := testutils.NewBuilding("building-one")
		floor := gateway.Floor{
			Building: building,
			Level:    1,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "for test",
			},
		}

		data, _ := json.Marshal(floor)

		actual, err := gateway.NewFloor(building, data)

		if assert.NoError(t, err) {
			if !cmp.Equal(floor, actual) {
				assert.Fail(t, cmp.Diff(floor, actual))
			}
		}
	})
}

func TestNewFloors(t *testing.T) {

	t.Run("should return floor associated to a building", func(t *testing.T) {
		building := testutils.NewBuilding("building-one")
		floors := gateway.Floors{"floor-one": gateway.Floor{
			Building: building,
			Level:    1,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "floor-one",
				Description: "for test",
			},
		}}

		data, _ := json.Marshal(floors)

		actual, err := gateway.NewFloors(building, data)

		if assert.NoError(t, err) {
			if !cmp.Equal(floors, actual) {
				assert.Fail(t, cmp.Diff(floors, actual))
			}
		}
	})
}
