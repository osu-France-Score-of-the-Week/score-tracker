package osu

type BeatmapsetResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Creator string `json:"creator"`
	Artist  string `json:"artist"`
	Status  string `json:"status"`
}
