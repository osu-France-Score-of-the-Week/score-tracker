package models

import "time"

type OAuthResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

type RecentScoresResponse struct {
	Scores []OsuScore `json:"scores"`
	Cursor string     `json:"cursor_string"`
}

type OsuScore struct {
	ID         uint           `json:"id"`
	Accuracy   float64        `json:"accuracy"`
	BeatmapID  uint           `json:"beatmap_id"`
	EndedAt    time.Time      `json:"ended_at"`
	HasReplay  bool           `json:"has_replay"`
	MaxCombo   uint           `json:"max_combo"`
	Mods       []Mod          `json:"mods"`
	Pp         *float64       `json:"pp"`
	Rank       string         `json:"rank"`
	Statistics map[string]int `json:"statistics"`
	UserID     uint           `json:"user_id"`
}

type Mod struct {
	Acronym string `json:"acronym"`
}

type ModSettings struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
