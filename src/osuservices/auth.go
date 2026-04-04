package osuservices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"score-tracker/models"
	"score-tracker/repositories"
	"strings"
	"time"
)

func (s *OsuService) getOsuToken() (models.OAuthResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	clientID := os.Getenv("OSU_CLIENT_ID")
	clientSecret := os.Getenv("OSU_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return models.OAuthResponse{}, fmt.Errorf("missing OSU_CLIENT_ID or OSU_CLIENT_SECRET environment variable")
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
		return models.OAuthResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return models.OAuthResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.OAuthResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return models.OAuthResponse{}, fmt.Errorf("osu OAuth error: %s", string(body))
	}

	var oauth models.OAuthResponse
	if err := json.Unmarshal(body, &oauth); err != nil {
		return models.OAuthResponse{}, err
	}

	fmt.Println("Successfully obtained osu access token")

	return oauth, nil
}

func (s *OsuService) getValidOsuToken() (string, error) {
	tokenRepo := repositories.NewTokenRepository(s.db)

	token, _ := tokenRepo.GetValidToken()

	if token == nil {
		fmt.Println("No valid token found, requesting new token from osu API")
		tokenStr, err := s.getOsuToken()
		if err != nil {
			return "", err
		}

		newToken := models.Token{
			Token:     tokenStr.AccessToken,
			ExpiresAt: time.Now().Add(time.Duration(tokenStr.ExpiresIn)*time.Second - time.Minute), // Subtract 1 minute to ensure we refresh before it actually expires
		}

		if err := tokenRepo.Create(&newToken); err != nil {
			return "", err
		}

		return tokenStr.AccessToken, nil
	}
	//fmt.Println("Valid token found in database, using existing token")
	return token.Token, nil

}
