package types

type VehicleType string
type EventType string

const (
	// vehicle types
	ShipVehicleType       VehicleType = "ship"
	HelicopterVehicleType VehicleType = "helicopter"

	// event types
	DepartureEventType EventType = "departed"
	ArrivalEventType   EventType = "arrived"
)

type VehicleEventMessage struct {
	VehicleType   VehicleType   `json:"vehicle_type"`
	VehicleUUID   UUID          `json:"vehicle_uuid"`
	StructureType StructureType `json:"structure_type"`
	StructureUUID UUID          `json:"structure_uuid"`
	Timestamp     int           `json:"timestamp"`
	Event         EventType     `json:"event"`
	SlotNumber    int           `json:"slot_number"`
	TowerUUID     UUID          `json:"tower_id"`
}
