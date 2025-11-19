package structure

import "github.com/google/uuid"

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
	Uuid uuid.UUID `json:"platform_uuid" db:"id"`
}

type Central struct {
	Structure
	Uuid uuid.UUID `json:"central_uuid" db:"id"`
}

type Structures struct {
	Platforms []Platform `json:"platforms"`
	Centrals  []Central  `json:"centrals"`
}
