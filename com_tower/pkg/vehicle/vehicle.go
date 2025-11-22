package vehicle

import "github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"

type VehicleType string

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
)

type VehicleEventMessage struct {
	VehicleType VehicleType `json:"vehicle_type"`
	VehicleUUID utils.UUID
}
