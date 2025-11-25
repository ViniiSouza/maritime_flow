package types

type Tower struct {
	UUID      UUID    `json:"tower_uuid" db:"id"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}

type TowerHealthRequest struct {
	Id UUID `json:"tower_id"`
}

type TowersPayload struct {
	Towers []Tower `json:"towers"`
}
