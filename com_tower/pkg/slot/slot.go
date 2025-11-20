package slot

import "github.com/ViniiSouza/maritime_flow/com_tower/pkg/structure"

type SlotType string
type SlotState string
type VehicleType string
type AcquireSlotResultType string

const (
	// slot types
	DockSlotType    SlotType = "dock"
	HelipadSlotType SlotType = "helipad"

	// slot states
	FreeSlotState  SlotState = "free"
	InUseSlotState SlotState = "in_use"

	// vehicle types
	ShipVehicleType       VehicleType = "ship"
	HelicopterVehicleType VehicleType = "helicopter"

	// acquire slot result types
	AcquiredAcquireSlotResultType    AcquireSlotResultType = "acquired"
	UnavailableAcquireSlotResultType AcquireSlotResultType = "unavailable"
)

type StructureSlotRequest struct {
	SlotNumber int      `json:"slot_number"`
	SlotType   SlotType `json:"slot_type"`
}

type SlotRequest struct {
	VehicleUuid   string                  `json:"vehicle_uuid"`
	VehicleType   VehicleType             `json:"vehicle_type"`
	StructureUuid string                  `json:"structure_uuid"`
	StructureType structure.StructureType `json:"structure_type"`
	StructureSlotRequest
}

type SlotResponse struct {
	State SlotState `json:"state"`
}

type AcquireSlotRequest struct {
	VehicleUuid   string `json:"vehicle_uuid"`
	StructureUuid string `json:"structure_uuid"`
	StructureSlotRequest
}

type AcquireSlotResponse struct {
	Result AcquireSlotResultType
}
