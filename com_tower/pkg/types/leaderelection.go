package types

type Role int

const (
	Leader Role = iota
	Minion
)

type ElectionRequest struct {
	CandidateUptime float64 `json:"candidate_uptime"`
}

type ElectionResponse struct {
	Uptime          float64 `json:"uptime"`
	HasHigherUptime bool    `json:"has_higher_uptime"`
}

type NewLeaderRequest struct {
	NewLeaderUUID UUID `json:"new_leader_uuid"`
}
