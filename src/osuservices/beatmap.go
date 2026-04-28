package osuservices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"score-tracker/models"
	osuModels "score-tracker/models/osu"
	"strings"
)

func (s *OsuService) GetBeatmap(beatmapId uint) (osuModels.BeatmapResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmaps/%d", beatmapId), nil)
	if err != nil {
		return osuModels.BeatmapResponse{}, err
	}

	token, err := s.getValidOsuToken()
	if err != nil {
		return osuModels.BeatmapResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return osuModels.BeatmapResponse{}, err
	}
	defer resp.Body.Close()

	var beatmap osuModels.BeatmapResponse
	if err := json.NewDecoder(resp.Body).Decode(&beatmap); err != nil {
		return osuModels.BeatmapResponse{}, err
	}

	return beatmap, nil
}

func (s *OsuService) GetBeatmapAttribute(beatmapId uint, mods models.Mods) (osuModels.BeatmapAttributesResponse, error) {
	client := &http.Client{}

	request := osuModels.BeatmapAttributesRequest{Mods: []osuModels.ModRequest{}}
	for _, mod := range mods {
		request.Mods = append(request.Mods, osuModels.ModRequest{
			Acronym:  mod.Acronym,
			Settings: mod.Settings,
		})
	}

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return osuModels.BeatmapAttributesResponse{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmaps/%d/attributes", beatmapId), strings.NewReader(string(bodyBytes)))
	if err != nil {
		return osuModels.BeatmapAttributesResponse{}, err
	}

	token, err := s.getValidOsuToken()
	if err != nil {
		return osuModels.BeatmapAttributesResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return osuModels.BeatmapAttributesResponse{}, err
	}
	defer resp.Body.Close()

	var attributes osuModels.BeatmapAttributesResponse
	if err := json.NewDecoder(resp.Body).Decode(&attributes); err != nil {
		return osuModels.BeatmapAttributesResponse{}, err
	}

	return attributes, nil
}
