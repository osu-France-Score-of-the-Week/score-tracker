package models

import "score-tracker/models/osu"

type Beatmap struct {
	ID               uint `gorm:"primaryKey"`
	BeatmapsetID     uint
	Beatmapset       Beatmapset
	DifficultyRating float64
	Version          string
}

func MapOsuBeatmapToModel(b osu.BeatmapResponse) Beatmap {
	return Beatmap{
		ID:               b.ID,
		BeatmapsetID:     b.Beatmapset.ID,
		Beatmapset:       MapOsuBeatmapsetToModel(b.Beatmapset),
		DifficultyRating: b.DifficultyRating,
		Version:          b.Version,
	}
}
