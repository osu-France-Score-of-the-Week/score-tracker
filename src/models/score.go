package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Score struct {
	ID         uint `gorm:"primaryKey"`
	Accuracy   float64
	BeatmapID  uint
	EndedAt    int64
	HasReplay  bool
	MaxCombo   uint
	Mods       Mods `gorm:"type:jsonb"`
	Pp         float64
	Rank       string
	Statistics string
	PlayerID   uint
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

func MapOsuScoreToModel(s OsuScore) (Score, error) {
	statsJSON, err := json.Marshal(s.Statistics)
	if err != nil {
		return Score{}, err
	}

	pp := 0.0
	if s.Pp != nil {
		pp = *s.Pp
	}

	mods := Mods(s.Mods)
	if mods == nil {
		mods = Mods{}
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
		Statistics: string(statsJSON),
		PlayerID:   s.UserID,
	}, nil
}

type Mod struct {
	Acronym  string                 `json:"acronym"`
	Settings map[string]interface{} `json:"settings"`
}
