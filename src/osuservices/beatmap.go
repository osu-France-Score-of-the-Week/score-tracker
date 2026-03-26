package osuservices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"score-tracker/models"
	"strings"
)

func GetBeatmap(beatmapId uint) (models.OsuBeatmap, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmaps/%d", beatmapId), nil)
	if err != nil {
		return models.OsuBeatmap{}, err
	}

	token, err := getOsuToken()
	if err != nil {
		return models.OsuBeatmap{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return models.OsuBeatmap{}, err
	}
	defer resp.Body.Close()

	var beatmap models.OsuBeatmap
	if err := json.NewDecoder(resp.Body).Decode(&beatmap); err != nil {
		return models.OsuBeatmap{}, err
	}

	return beatmap, nil
}

func GetBeatmapAttribute(beatmapId uint, mods models.Mods) (models.OsuBeatmapAttributes, error) {
	client := &http.Client{}

	var modsVal interface{}
	if mods != nil && len(mods) > 0 {
		modsVal = mods
	} else {
		modsVal = []interface{}{}
	}

	payload := map[string]interface{}{"mods": modsVal}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return models.OsuBeatmapAttributes{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmaps/%d/attributes", beatmapId), strings.NewReader(string(bodyBytes)))
	if err != nil {
		return models.OsuBeatmapAttributes{}, err
	}

	token, err := getOsuToken()
	if err != nil {
		return models.OsuBeatmapAttributes{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.OsuBeatmapAttributes{}, err
	}
	defer resp.Body.Close()

	var attributes models.OsuBeatmapAttributes
	if err := json.NewDecoder(resp.Body).Decode(&attributes); err != nil {
		return models.OsuBeatmapAttributes{}, err
	}

	return attributes, nil
}
