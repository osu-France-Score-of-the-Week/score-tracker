package responses

type ScoreCursor struct {
	Scores []ScoreWithAttributes `json:"scores"`
	Cursor string                `json:"cursor"`
}
