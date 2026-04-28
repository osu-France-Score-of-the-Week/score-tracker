package responses

type ScorePage struct {
	Scores     []ScoreWithAttributes `json:"scores"`
	Page       int                   `json:"page"`
	TotalPages int                   `json:"totalPages"`
}
