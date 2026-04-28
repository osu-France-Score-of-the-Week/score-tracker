package osuservices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"score-tracker/models"
	osuModels "score-tracker/models/osu"
	"score-tracker/repositories"
	"strings"
	"time"
)

func (s *OsuService) getOsuToken() (osuModels.OAuthResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	clientID := os.Getenv("OSU_CLIENT_ID")
	clientSecret := os.Getenv("OSU_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return osuModels.OAuthResponse{}, fmt.Errorf("missing OSU_CLIENT_ID or OSU_CLIENT_SECRET environment variable")
	}

	request := osuModels.OAuthRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
		Scope:        "public",
	}

	form := request.ToValues()

	req, err := http.NewRequest(
		"POST",
		"https://osu.ppy.sh/oauth/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return osuModels.OAuthResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return osuModels.OAuthResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return osuModels.OAuthResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return osuModels.OAuthResponse{}, fmt.Errorf("osu OAuth error: %s", string(body))
	}

	var oauth osuModels.OAuthResponse
	if err := json.Unmarshal(body, &oauth); err != nil {
		return osuModels.OAuthResponse{}, err
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
