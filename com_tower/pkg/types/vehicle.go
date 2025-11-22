package types

import (
	"time"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

type VehicleType string
type EventType string

/*
H/S -> RabbitMQ (Audit | Events)
{
  "vehicle_type": "helicopter"|"ship",
  "vehicle_uuid": "823429efabc9283",
  "structure_type": "platform"|"central",
  "structure_uuid": "acb432efab98234d",
  "timestamp": 3094870293,
  "event": "arrived"|"departed",
  "slot_number": 0,
  "tower_id": "adbae9438ff92a"
}
*/

const (
	// vehicle types
	ShipVehicleType       VehicleType = "ship"
	HelicopterVehicleType VehicleType = "helicopter"

	// event types
	DepartureEventType EventType = "departed"
	ArrivalEventType   EventType = "arrived"
)

type VehicleEventMessage struct {
	VehicleType   VehicleType             `json:"vehicle_type"`
	VehicleUUID   utils.UUID              `json:"vehicle_uuid"`
	StructureType StructureType `json:"structure_type"`
	StructureUUID utils.UUID              `json:"structure_uuid"`
	Timestamp     time.Time               `json:"timestamp"`
	Event         EventType               `json:"event"`
	SlotNumber    int                     `json:"slot_number"`
	TowerUUID     utils.UUID              `json:"tower_id"`
}
