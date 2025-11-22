package tower

import (
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

type Tower struct {
	UUID      utils.UUID `json:"tower_uuid" db:"id"`
	Latitude  float64    `json:"latitude" db:"latitude"`
	Longitude float64    `json:"longitude" db:"longitude"`
}

type TowerHealthRequest struct {
	Id utils.UUID `json:"tower_id"`
}

type TowersPayload struct {
	Towers []Tower `json:"towers"`
}
