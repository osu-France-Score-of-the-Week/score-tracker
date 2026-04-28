package models

import "score-tracker/models/osu"

type Beatmapset struct {
	ID      uint `gorm:"primaryKey"`
	Artist  string
	Creator string
	Title   string
	Status  string
}

func MapOsuBeatmapsetToModel(b osu.BeatmapsetResponse) Beatmapset {
	return Beatmapset{
		ID:      b.ID,
		Artist:  b.Artist,
		Creator: b.Creator,
		Title:   b.Title,
		Status:  b.Status,
	}
}
