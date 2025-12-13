package osuservices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"score-tracker/models"
	"strings"
	"time"
)

func getOsuToken() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	form := url.Values{}
	form.Set("client_id", "25310")
	form.Set("client_secret", "8KcttqiVcSyY6fvGBWxwqm52fe8EziV2oUsrOzdB")
	form.Set("grant_type", "client_credentials")
	form.Set("scope", "public")

	req, err := http.NewRequest(
		"POST",
		"https://osu.ppy.sh/oauth/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("osu OAuth error: %s", string(body))
	}

	var oauth models.OAuthResponse
	if err := json.Unmarshal(body, &oauth); err != nil {
		return "", err
	}

	return oauth.AccessToken, nil
}
