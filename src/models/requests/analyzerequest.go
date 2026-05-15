package requests

import "score-tracker/models"

type AnalyzeRequest struct {
	Beatmap           models.Beatmap           `json:"beatmap"`
	BeatmapAttributes models.BeatmapAttributes `json:"beatmap_attributes"`
	Score             models.Score             `json:"score"`
}
