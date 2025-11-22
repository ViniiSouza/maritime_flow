package types

type SlotType string
type SlotState string
type AcquireSlotResultType string

const (
	// slot types
	DockSlotType    SlotType = "dock"
	HelipadSlotType SlotType = "helipad"

	// slot states
	FreeSlotState  SlotState = "free"
	InUseSlotState SlotState = "in_use"

	// acquire slot result types
	AcquiredAcquireSlotResultType    AcquireSlotResultType = "acquired"
	UnavailableAcquireSlotResultType AcquireSlotResultType = "unavailable"
)

type StructureSlotRequest struct {
	SlotNumber int      `json:"slot_number"`
	SlotType   SlotType `json:"slot_type"`
}

type SlotRequest struct {
	VehicleUUID   UUID              `json:"vehicle_uuid"`
	VehicleType   VehicleType     `json:"vehicle_type"`
	StructureUUID UUID              `json:"structure_uuid"`
	StructureType StructureType `json:"structure_type"`
	StructureSlotRequest
}

type SlotResponse struct {
	State SlotState `json:"state"`
}

type AcquireSlotRequest struct {
	VehicleUUID   UUID `json:"vehicle_uuid"`
	StructureUUID UUID `json:"structure_uuid"`
	StructureSlotRequest
}

type AcquireSlotResponse struct {
	Result AcquireSlotResultType
}

type ReleaseSlotRequest StructureSlotRequest

type ReleaseSlotLockRequest AcquireSlotRequest
