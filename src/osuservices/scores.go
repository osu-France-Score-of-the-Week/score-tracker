package osuservices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"score-tracker/models"
)

func GetRecentScores(cursor *string) (models.RecentScoresResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://osu.ppy.sh/api/v2/scores", nil)
	if err != nil {
		return models.RecentScoresResponse{}, err
	}

	if cursor != nil {
		q := req.URL.Query()
		q.Add("cursor_string", *cursor)
		req.URL.RawQuery = q.Encode()
	}

	token, err := getOsuToken()
	if err != nil {
		return models.RecentScoresResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return models.RecentScoresResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.RecentScoresResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return models.RecentScoresResponse{}, fmt.Errorf("osu API error: %s", string(body))
	}

	var recentScores models.RecentScoresResponse
	if err := json.Unmarshal(body, &recentScores); err != nil {
		return models.RecentScoresResponse{}, err
	}

	fmt.Println("Fetched", len(recentScores.Scores), "recent scores from osu API")

	return recentScores, nil
}
