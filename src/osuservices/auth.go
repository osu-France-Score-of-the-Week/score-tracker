package osuservices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"score-tracker/models"
	"strings"
	"time"
)

func getOsuToken() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	clientID := os.Getenv("OSU_CLIENT_ID")
	clientSecret := os.Getenv("OSU_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("missing OSU_CLIENT_ID or OSU_CLIENT_SECRET environment variable")
	}

	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

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
