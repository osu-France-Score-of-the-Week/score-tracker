package osu

type BeatmapAttributesResponse struct {
	Attributes AttributesResponse `json:"attributes"`
}

type AttributesResponse struct {
	StarRating float64 `json:"star_rating"`
	MaxCombo   uint    `json:"max_combo"`
}
