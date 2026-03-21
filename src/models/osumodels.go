package models

import (
	"encoding/json"
	"time"
)

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

type UsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	ID         uint       `json:"id"`
	Username   string     `json:"username"`
	Country    string     `json:"country_code"`
	Statistics Statistics `json:"statistics_rulesets"`
}

type Statistics struct {
	GlobalRank int     `json:"global_rank"`
	Pp         float64 `json:"pp"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID         uint                  `json:"id"`
		Username   string                `json:"username"`
		Country    string                `json:"country_code"`
		Statistics map[string]Statistics `json:"statistics_rulesets"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	u.ID = raw.ID
	u.Username = raw.Username
	u.Country = raw.Country

	if stats, ok := raw.Statistics["osu"]; ok {
		u.Statistics = stats
	}

	return nil
}
