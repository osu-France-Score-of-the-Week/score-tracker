package models

import "score-tracker/models/osu"

type Beatmap struct {
	ID               uint `gorm:"primaryKey"`
	BeatmapsetID     uint
	Beatmapset       Beatmapset
	DifficultyRating float64
	Status           string  `json:"status"`
	CS               float64 `json:"cs"`
	OD               float64 `json:"accuracy"`
	AR               float64 `json:"ar"`
	Version          string
}

func MapOsuBeatmapToModel(b osu.BeatmapResponse) Beatmap {
	return Beatmap{
		ID:               b.ID,
		BeatmapsetID:     b.Beatmapset.ID,
		Beatmapset:       MapOsuBeatmapsetToModel(b.Beatmapset),
		DifficultyRating: b.DifficultyRating,
		Status:           b.Status,
		CS:               b.CS,
		OD:               b.OD,
		AR:               b.AR,
		Version:          b.Version,
	}
}
