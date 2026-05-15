package models

import (
	"database/sql/driver"
	"encoding/json"
	"score-tracker/models/osu"
)

type Score struct {
	ID            uint `gorm:"primaryKey"`
	Accuracy      float64
	BeatmapID     uint
	Beatmap       Beatmap
	EndedAt       int64
	HasReplay     bool
	MaxCombo      uint
	Mods          Mods `gorm:"type:jsonb"`
	Pp            float64
	Rank          string
	Statistics    ScoreStatistics `gorm:"type:jsonb"`
	PlayerID      uint
	Player        Player
	AnalyzedScore float64
}

type ScoreStatistics map[string]int

func (s ScoreStatistics) Value() (driver.Value, error) {
	if s == nil {
		return []byte("{}"), nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *ScoreStatistics) Scan(value interface{}) error {
	if value == nil {
		*s = ScoreStatistics{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return nil
	}

	if len(bytes) == 0 {
		*s = ScoreStatistics{}
		return nil
	}

	var out map[string]int
	if err := json.Unmarshal(bytes, &out); err != nil {
		return err
	}
	*s = ScoreStatistics(out)
	return nil
}

type Mods []Mod

func (m Mods) Value() (driver.Value, error) {
	if m == nil {
		return []byte("null"), nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (m *Mods) Scan(value interface{}) error {
	if value == nil {
		*m = Mods{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return nil
	}

	if len(bytes) == 0 {
		*m = Mods{}
		return nil
	}

	var out []Mod
	if err := json.Unmarshal(bytes, &out); err != nil {
		return err
	}
	*m = Mods(out)
	return nil
}

func MapOsuScoreToModel(s osu.ScoreResponse) (Score, error) {
	pp := 0.0
	if s.Pp != nil {
		pp = *s.Pp
	}

	mods := make(Mods, 0, len(s.Mods))
	for _, mod := range s.Mods {
		mods = append(mods, Mod{Acronym: mod.Acronym, Settings: mod.Settings})
	}
	if mods == nil {
		mods = Mods{}
	}

	stats := ScoreStatistics(s.Statistics)
	if stats == nil {
		stats = ScoreStatistics{}
	}

	return Score{
		ID:         s.ID,
		Accuracy:   s.Accuracy,
		BeatmapID:  s.BeatmapID,
		EndedAt:    s.EndedAt.Unix(),
		HasReplay:  s.HasReplay,
		MaxCombo:   s.MaxCombo,
		Mods:       mods,
		Pp:         pp,
		Rank:       s.Rank,
		Statistics: stats,
		PlayerID:   s.UserID,
	}, nil
}

type Mod struct {
	Acronym  string                 `json:"acronym"`
	Settings map[string]interface{} `json:"settings"`
}
