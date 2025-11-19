package tower

import "github.com/google/uuid"

type Tower struct {
	Uuid      uuid.UUID `json:"tower_uuid" db:"id"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
}

type TowerHealthRequest struct {
	Id string `json:"tower_id"`
}

type TowersResponse struct {
	Towers []Tower `json:"towers"`
}
