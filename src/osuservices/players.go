package osuservices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"score-tracker/models"
)

func GetPlayers(playerIds []uint) (models.UsersResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://osu.ppy.sh/api/v2/users", nil)
	if err != nil {
		return models.UsersResponse{}, err
	}

	for _, playerId := range playerIds {
		q := req.URL.Query()
		q.Add("ids[]", fmt.Sprintf("%d", playerId))
		req.URL.RawQuery = q.Encode()
	}

	token, err := getOsuToken()
	if err != nil {
		return models.UsersResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return models.UsersResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return models.UsersResponse{}, fmt.Errorf("osu API error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.UsersResponse{}, err
	}

	var players models.UsersResponse
	if err := json.Unmarshal(body, &players); err != nil {
		return models.UsersResponse{}, err
	}

	return players, nil
}
