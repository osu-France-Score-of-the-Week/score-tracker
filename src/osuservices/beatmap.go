package osuservices

import (
	"encoding/json"
	"fmt"
	"score-tracker/models"
	osuModels "score-tracker/models/osu"
	"score-tracker/queue"
)

func (s *OsuService) GetBeatmap(beatmapId uint) (osuModels.BeatmapResponse, error) {
	token, err := s.getValidOsuToken()
	if err != nil {
		return osuModels.BeatmapResponse{}, err
	}

	return queue.ExecuteJSON[osuModels.BeatmapResponse](s.queue, queue.Request{
		Method: "GET",
		URL:    fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmaps/%d", beatmapId),
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
}

func (s *OsuService) GetBeatmapAttribute(beatmapId uint, mods models.Mods) (osuModels.BeatmapAttributesResponse, error) {
	token, err := s.getValidOsuToken()
	if err != nil {
		return osuModels.BeatmapAttributesResponse{}, err
	}

	request := osuModels.BeatmapAttributesRequest{Mods: []osuModels.ModRequest{}}
	for _, mod := range mods {
		request.Mods = append(request.Mods, osuModels.ModRequest{
			Acronym:  mod.Acronym,
			Settings: mod.Settings,
		})
	}

	body, err := json.Marshal(request)
	if err != nil {
		return osuModels.BeatmapAttributesResponse{}, err
	}

	return queue.ExecuteJSON[osuModels.BeatmapAttributesResponse](s.queue, queue.Request{
		Method: "POST",
		URL:    fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmaps/%d/attributes", beatmapId),
		Body:   body,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
		},
	})
}
