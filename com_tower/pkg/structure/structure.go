package structure

import (
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

type StructureType string

const (
	PlatformStructureType StructureType = "platform"
	CentralStructureType  StructureType = "central"
)

type slots struct {
	DocksQtt    int `json:"docks_qtt"`
	HelipadsQtt int `json:"helipads_qtt"`
}

type Structure struct {
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
	Slots     slots   `json:"slots"`
}

type Platform struct {
	Structure
	UUID utils.UUID `json:"platform_uuid" db:"id"`
}

type Central struct {
	Structure
	UUID utils.UUID `json:"central_uuid" db:"id"`
}

type Structures struct {
	Platforms []Platform `json:"platforms"`
	Centrals  []Central  `json:"centrals"`
}
