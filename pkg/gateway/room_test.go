package gateway_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/testutils"
	"testing"
)

func TestNewRoom(t *testing.T) {
	t.Run("should return room associated to a room", func(t *testing.T) {
		floor := testutils.NewFloor("floor-one")
		room := gateway.Room{
			Floor:     floor,
			Direction: gateway.DirectionEast,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "for test",
			},
		}

		data, _ := json.Marshal(room)

		actual, err := gateway.NewRoom(floor, data)

		if assert.NoError(t, err) {
			if !cmp.Equal(room, actual) {
				assert.Fail(t, cmp.Diff(room, actual))
			}
		}
	})
}

func TestNewRooms(t *testing.T) {

	t.Run("should return room associated to a floor", func(t *testing.T) {
		floor := testutils.NewFloor("floor-one")
		rooms := gateway.Rooms{"room-one": gateway.Room{
			Floor:     floor,
			Direction: gateway.DirectionEast,
			PhysicalEntity: gateway.PhysicalEntity{
				Name:        "room-one",
				Description: "for test",
			},
		}}

		data, _ := json.Marshal(rooms)

		actual, err := gateway.NewRooms(floor, data)

		if assert.NoError(t, err) {
			if !cmp.Equal(rooms, actual) {
				assert.Fail(t, cmp.Diff(rooms, actual))
			}
		}
	})
}
