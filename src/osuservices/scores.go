package osuservices

import (
	"fmt"
	"net/url"
	osuModels "score-tracker/models/osu"
	"score-tracker/queue"
)

func (s *OsuService) GetRecentScores(cursor string) (osuModels.RecentScoresResponse, error) {
	token, err := s.getValidOsuToken()
	if err != nil {
		return osuModels.RecentScoresResponse{}, err
	}

	q := url.Values{}
	q.Set("ruleset", "osu")
	if cursor != "" {
		q.Set("cursor_string", cursor)
	}

	recentScores, err := queue.ExecuteJSON[osuModels.RecentScoresResponse](s.queue, queue.Request{
		Method: "GET",
		URL:    "https://osu.ppy.sh/api/v2/scores?" + q.Encode(),
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	if err != nil {
		return osuModels.RecentScoresResponse{}, err
	}

	fmt.Println("Fetched", len(recentScores.Scores), "recent scores from osu API")

	return recentScores, nil
}
