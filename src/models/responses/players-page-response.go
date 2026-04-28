package responses

type PlayersPage struct {
	Players    []Player `json:"players"`
	Page       int      `json:"page"`
	TotalPages int64    `json:"totalPages"`
}
