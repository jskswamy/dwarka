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
	}

	for _, testScenario := range scenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			actual := gateway.NewDirection(testScenario.direction)

			assert.Equal(t, actual, testScenario.expected)
		})
	}
}
