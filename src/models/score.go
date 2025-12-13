package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Score struct {
	gorm.Model
	ScoreId    uint `gorm:"primaryKey"`
	Accuracy   float64
	BeatmapId  uint
	EndedAt    int64
	HasReplay  bool
	MaxCombo   uint
	Mods       string
	Pp         float64
	Rank       string
	Statistics string
	PlayerId   uint
}

func MapOsuScoreToModel(s OsuScore) (Score, error) {
	modsJSON, err := json.Marshal(s.Mods)
	if err != nil {
		return Score{}, err
	}

	statsJSON, err := json.Marshal(s.Statistics)
	if err != nil {
		return Score{}, err
	}

	pp := 0.0
	if s.Pp != nil {
		pp = *s.Pp
	}

	return Score{
		ScoreId:    s.ID,
		Accuracy:   s.Accuracy,
		BeatmapId:  s.BeatmapID,
		EndedAt:    s.EndedAt.Unix(),
		HasReplay:  s.HasReplay,
		MaxCombo:   s.MaxCombo,
		Mods:       string(modsJSON),
		Pp:         pp,
		Rank:       s.Rank,
		Statistics: string(statsJSON),
		PlayerId:   s.UserID,
	}, nil
}
