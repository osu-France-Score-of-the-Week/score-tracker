package osu

import "encoding/json"

type UserResponse struct {
	ID         uint               `json:"id"`
	Username   string             `json:"username"`
	Country    string             `json:"country_code"`
	Statistics StatisticsResponse `json:"statistics_rulesets"`
}

func (u *UserResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID         uint                          `json:"id"`
		Username   string                        `json:"username"`
		Country    string                        `json:"country_code"`
		Statistics map[string]StatisticsResponse `json:"statistics_rulesets"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	u.ID = raw.ID
	u.Username = raw.Username
	u.Country = raw.Country

	if stats, ok := raw.Statistics["osu"]; ok {
		u.Statistics = stats
	}

	return nil
}
