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
		if applyFilters(player, []PlayerFilter{isFrench, isTop10k}) {
			fmt.Println("Player passed filters:", player.Username)
			fmt.Println("statistics:", player.Statistics)
			for _, score := range scores {
				if score.PlayerId == player.ID {
					scoresChan <- score
				}
			}
		}
	}
}

type PlayerFilter func(player models.User) bool

func isFrench(player models.User) bool {
	return player.Country == "FR"
}

func isTop10k(player models.User) bool {
	return player.Statistics.GlobalRank <= 10000 && player.Statistics.GlobalRank > 0
}

func applyFilters(player models.User, filters []PlayerFilter) bool {
	for _, filter := range filters {
		if !filter(player) {
			return false
		}
	}
	return true
}
