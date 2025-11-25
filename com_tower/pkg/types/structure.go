package types

type StructureType string

const (
	PlatformStructureType StructureType = "platform"
	CentralStructureType  StructureType = "central"
)

type StructureSlots struct {
	DocksQtt    int `json:"docks_qtt" db:"docks_qtt"`
	HelipadsQtt int `json:"helipads_qtt" db:"helipads_qtt"`
}

type Structure struct {
	Latitude  float64        `json:"latitude" db:"latitude"`
	Longitude float64        `json:"longitude" db:"longitude"`
	Slots     StructureSlots `json:"slots" db:"slots"`
}

type Platform struct {
	Structure
	UUID UUID `json:"platform_uuid" db:"id"`
}

type Central struct {
	Structure
	UUID UUID `json:"central_uuid" db:"id"`
}

type Structures struct {
	Platforms []Platform `json:"platforms"`
	Centrals  []Central  `json:"centrals"`
}
