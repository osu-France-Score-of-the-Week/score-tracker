package jobs

import (
	"score-tracker/models"
	"score-tracker/osuservices"
	"score-tracker/repositories"

	"gorm.io/gorm"
)

func FilterScores(stopChan <-chan struct{}, filterChan <-chan models.OsuScore, scoresChan chan<- models.Score, db *gorm.DB) {
	playerRepo := repositories.NewPlayerRepository(db)

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
					checkPlayersFromScores(batch, scoresChan, playerRepo)
					batch = batch[:0] // Clear the batch
				}

			case <-stopChan:
				return
			}
		}
	}()
}

func checkPlayersFromScores(scores []models.Score, scoresChan chan<- models.Score, playerRepo *repositories.PlayerRepository) {

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

	for _, player := range players.Users {
		if applyFilters(player, []PlayerFilter{isFrench, isTop10k}) {
			newPlayer := models.MapOsuUserToModel(player)
			err := playerRepo.Update(&newPlayer)
			if err != nil {
				return
			}

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
