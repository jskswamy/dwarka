package gateway_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"testing"
)

func TestDirection_Direction(t *testing.T) {
	type scenario struct {
		name      string
		direction gateway.Direction
		expected  string
	}
	scenarios := []scenario{
		{
			name:      "Direction should convert north",
			direction: gateway.DirectionNorth,
			expected:  "north",
		},
		{
			name:      "Direction should convert east",
			direction: gateway.DirectionEast,
			expected:  "east",
		},
		{
			name:      "Direction should convert south",
			direction: gateway.DirectionSouth,
			expected:  "south",
		},
		{
			name:      "Direction should convert west",
			direction: gateway.DirectionWest,
			expected:  "west",
		},
	}

	for _, testScenario := range scenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			actual := testScenario.direction.Direction()

			assert.Equal(t, actual, testScenario.expected)
		})
	}
}

func TestNewDirection(t *testing.T) {
	type scenario struct {
		name      string
		direction string
		expected  gateway.Direction
		error     bool
	}
	scenarios := []scenario{
		{
			name:      "NewDirection should convert north",
			expected:  gateway.DirectionNorth,
			direction: "north",
		},
		{
			name:      "NewDirection should convert east",
			expected:  gateway.DirectionEast,
			direction: "east",
		},
		{
			name:      "NewDirection should convert south",
			expected:  gateway.DirectionSouth,
			direction: "south",
		},
		{
			name:      "NewDirection should convert west",
			expected:  gateway.DirectionWest,
			direction: "west",
		},
		{
			name:      "NewDirection should not convert south west",
			expected:  -1,
			direction: "south west",
			error:     true,
		},
	}

	for _, testScenario := range scenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			actual, err := gateway.NewDirection(testScenario.direction)

			assert.Equal(t, actual, testScenario.expected)
			if testScenario.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
