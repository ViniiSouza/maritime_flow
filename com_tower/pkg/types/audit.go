package types

type ResultType string

const (
	// result types
	AllowedResultType ResultType = "allowed"
	DeniedResultType  ResultType = "denied"
)

type AuditRequest struct {
	VehicleType   VehicleType   `json:"vehicle_type"`
	VehicleUUID   UUID          `json:"vehicle_uuid"`
	StructureType StructureType `json:"structure_type"`
	StructureUUID UUID          `json:"structure_uuid"`
	Timestamp     int           `json:"timestamp"`
	Result        ResultType    `json:"result"`
	SlotNumber    int           `json:"slot_number"`
	TowerUUID     UUID          `json:"tower_id"`
}
