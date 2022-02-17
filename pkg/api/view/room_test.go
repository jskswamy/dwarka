package view_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api/view"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"testing"
)

func TestNewRoom(t *testing.T) {
	type scenario struct {
		name     string
		room     gateway.Room
		expected view.Room
	}
	scenarios := []scenario{
		{
			name: "NewRoom should convert north room",
			expected: view.Room{
				Direction:   "north",
				Name:        "north room",
				Description: "room facing north",
			},
			room: gateway.Room{
				Direction: gateway.DirectionNorth,
				PhysicalEntity: gateway.PhysicalEntity{
					Name:        "north room",
					Description: "room facing north",
				},
			},
		},
		{
			name: "NewRoom should convert east room",
			expected: view.Room{
				Direction:   "east",
				Name:        "east room",
				Description: "room facing east",
			},
			room: gateway.Room{
				Direction: gateway.DirectionEast,
				PhysicalEntity: gateway.PhysicalEntity{
					Name:        "east room",
					Description: "room facing east",
				},
			},
		},
		{
			name: "NewRoom should convert south room",
			expected: view.Room{
				Direction:   "south",
				Name:        "south room",
				Description: "room facing south",
			},
			room: gateway.Room{
				Direction: gateway.DirectionSouth,
				PhysicalEntity: gateway.PhysicalEntity{
					Name:        "south room",
					Description: "room facing south",
				},
			},
		},
		{
			name: "NewRoom should convert west room",
			expected: view.Room{
				Direction:   "west",
				Name:        "west room",
				Description: "room facing west",
			},
			room: gateway.Room{
				Direction: gateway.DirectionWest,
				PhysicalEntity: gateway.PhysicalEntity{
					Name:        "west room",
					Description: "room facing west",
				},
			},
		},
	}

	for _, testScenario := range scenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			actual := view.NewRoom(testScenario.room)

			if !cmp.Equal(testScenario.expected, actual) {
				assert.Fail(t, cmp.Diff(testScenario.expected, actual))
			}
		})
	}
}

func TestRoom_Room(t *testing.T) {
	type scenario struct {
		name     string
		room     view.Room
		expected gateway.Room
	}
	scenarios := []scenario{
		{
			name: "Room should convert north room",
			room: view.Room{
				Direction:   "north",
				Name:        "north room",
				Description: "room facing north",
			},
			expected: gateway.Room{
				Direction: gateway.DirectionNorth,
				PhysicalEntity: gateway.PhysicalEntity{
					Name:        "north room",
					Description: "room facing north",
				},
			},
		},
		{
			name: "Room should convert east room",
			room: view.Room{
				Direction:   "east",
				Name:        "east room",
				Description: "room facing east",
			},
			expected: gateway.Room{
				Direction: gateway.DirectionEast,
				PhysicalEntity: gateway.PhysicalEntity{
					Name:        "east room",
					Description: "room facing east",
				},
			},
		},
		{
			name: "Room should convert south room",
			room: view.Room{
				Direction:   "south",
				Name:        "south room",
				Description: "room facing south",
			},
			expected: gateway.Room{
				Direction: gateway.DirectionSouth,
				PhysicalEntity: gateway.PhysicalEntity{
					Name:        "south room",
					Description: "room facing south",
				},
			},
		},
		{
			name: "Room should convert west room",
			room: view.Room{
				Direction:   "west",
				Name:        "west room",
				Description: "room facing west",
			},
			expected: gateway.Room{
				Direction: gateway.DirectionWest,
				PhysicalEntity: gateway.PhysicalEntity{
					Name:        "west room",
					Description: "room facing west",
				},
			},
		},
	}

	for _, testScenario := range scenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			actual := testScenario.room.Room()

			if !cmp.Equal(testScenario.expected, actual) {
				assert.Fail(t, cmp.Diff(testScenario.expected, actual))
			}
		})
	}
}
