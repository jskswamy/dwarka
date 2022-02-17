package gateway

import (
	"github.com/gosimple/slug"
)

const (
	// DirectionNorth the compass point corresponding to north
	DirectionNorth Direction = iota + 1

	// DirectionEast the direction towards the point of the horizon where the sun rises at the equinoxes
	DirectionEast

	// DirectionSouth the direction towards the point of the horizon 90° clockwise from east
	DirectionSouth

	// DirectionWest the direction towards the point of the horizon where the sun sets at the equinoxes,
	DirectionWest
)

// NodeType represents the node type
type NodeType int

const (
	// NodeTypeWifi represents the wifi node type
	NodeTypeWifi NodeType = iota

	// NodeTypeMqtt represents the mqtt node type
	NodeTypeMqtt
)

// Entity represents a uniquely identifiable object
type Entity interface {
	ID() string
	Validate() error
}

// PhysicalEntity a thing with distinct and independent existence
type PhysicalEntity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ID returns slug representing the entity
func (entity PhysicalEntity) ID() string {
	return slug.Make(entity.Name)
}

// Device is a electronic / electrical equipment made or adapted for a particular purpose
type Device struct {
	Meta           map[string]string `json:"meta"`
	PhysicalEntity `json:"entity"`
}

// Devices collection of device
type Devices []Device

// NodeMetadata represents information about a node
type NodeMetadata struct {
	Building `json:"building"`
	Floor    `json:"floor"`
	Room     `json:"room"`
	Devices  `json:"devices"`
	Host     string   `json:"host"`
	Type     NodeType `json:"type"`
}

// Node a piece of equipment
type Node interface {
	On(Device) error
	Off(Device) error
}

// Status represents key value collection of various status
type Status map[string]string
