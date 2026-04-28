package osu

import "time"

type ScoreResponse struct {
	ID         uint           `json:"id"`
	Accuracy   float64        `json:"accuracy"`
	BeatmapID  uint           `json:"beatmap_id"`
	EndedAt    time.Time      `json:"ended_at"`
	HasReplay  bool           `json:"has_replay"`
	MaxCombo   uint           `json:"max_combo"`
	Mods       []ModResponse  `json:"mods"`
	Pp         *float64       `json:"pp"`
	Rank       string         `json:"rank"`
	Statistics map[string]int `json:"statistics"`
	UserID     uint           `json:"user_id"`
}
