package jobs

import (
	"fmt"
	"score-tracker/models"
	"score-tracker/osuservices"
)

func FilterScores(stopChan <-chan struct{}, filterChan <-chan models.OsuScore, scoresChan chan<- models.Score) {
	batch := make([]models.Score, 0, 50)

	go func() {
		for {
			select {
			case scoreToFilter := <-filterChan:
				mappedScore, err := models.MapOsuScoreToModel(scoreToFilter)
				if err != nil {
					continue
				}

				batch = append(batch, mappedScore)

				if len(batch) == 50 {
					//scoresChan <- mappedScore
					checkPlayersFromScores(batch, scoresChan)
					batch = batch[:0] // Clear the batch
				}

			case <-stopChan:
				return
			}
		}
	}()
}

func checkPlayersFromScores(scores []models.Score, scoresChan chan<- models.Score) {
	playerIds := make([]uint, 0, len(scores))
	for _, score := range scores {
		playerIds = append(playerIds, score.PlayerId)
	}

	var players, err = osuservices.GetPlayers(playerIds)
	if err != nil {
		return
	}

	if len(players.Users) == 0 {
		return
	}

	fmt.Printf("Found %d players\n", len(players.Users))
	for _, player := range players.Users {
		if player.Country == "FR" {
			for _, score := range scores {
				if score.PlayerId == player.ID {
					scoresChan <- score
				}
			}
		}
	}
}
