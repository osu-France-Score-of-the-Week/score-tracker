package osu

type RecentScoresResponse struct {
	Scores []ScoreResponse `json:"scores"`
	Cursor string          `json:"cursor_string"`
}
