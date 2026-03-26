package models

type Beatmap struct {
	ID               uint `gorm:"primaryKey"`
	BeatmapsetID     uint
	Beatmapset       Beatmapset
	DifficultyRating float64
	Version          string
}

func MapOsuBeatmapToModel(b OsuBeatmap) Beatmap {
	return Beatmap{
		ID:               b.ID,
		BeatmapsetID:     b.Beatmapset.ID,
		Beatmapset:       MapOsuBeatmapsetToModel(b.Beatmapset),
		DifficultyRating: b.DifficultyRating,
		Version:          b.Version,
	}
}
