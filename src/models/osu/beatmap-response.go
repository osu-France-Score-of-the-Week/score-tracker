package osu

type BeatmapResponse struct {
	ID               uint               `json:"id"`
	Version          string             `json:"version"`
	DifficultyRating float64            `json:"difficulty_rating"`
	Beatmapset       BeatmapsetResponse `json:"beatmapset"`
}
