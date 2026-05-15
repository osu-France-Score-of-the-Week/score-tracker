package osu

type BeatmapResponse struct {
	ID               uint               `json:"id"`
	Version          string             `json:"version"`
	DifficultyRating float64            `json:"difficulty_rating"`
	Status           string             `json:"status"`
	CS               float64            `json:"cs"`
	OD               float64            `json:"accuracy"`
	AR               float64            `json:"ar"`
	Beatmapset       BeatmapsetResponse `json:"beatmapset"`
}
