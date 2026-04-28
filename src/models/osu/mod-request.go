package osu

type ModRequest struct {
	Acronym  string         `json:"acronym"`
	Settings map[string]any `json:"settings"`
}
