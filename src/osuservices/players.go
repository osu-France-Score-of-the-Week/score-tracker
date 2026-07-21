package osuservices

import (
	"fmt"
	"net/url"
	osuModels "score-tracker/models/osu"
	"score-tracker/queue"
)

func (s *OsuService) GetPlayers(request osuModels.UsersRequest) (osuModels.UsersResponse, error) {
	token, err := s.getValidOsuToken()
	if err != nil {
		return osuModels.UsersResponse{}, err
	}

	q := url.Values{}
	for _, playerId := range request.IDs {
		q.Add("ids[]", fmt.Sprintf("%d", playerId))
	}

	return queue.ExecuteJSON[osuModels.UsersResponse](s.queue, queue.Request{
		Method: "GET",
		URL:    "https://osu.ppy.sh/api/v2/users?" + q.Encode(),
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
}
