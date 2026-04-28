package osu

type ModResponse struct {
	Acronym  string         `json:"acronym"`
	Settings map[string]any `json:"settings"`
}
