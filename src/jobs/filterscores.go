package jobs

import (
	"fmt"
	"score-tracker/database"
	"score-tracker/models"
	osuModels "score-tracker/models/osu"
	"score-tracker/osuservices"
	"score-tracker/repositories"
	"time"

	"gorm.io/gorm"
)

func FilterScores(stopChan <-chan struct{}, filterChan <-chan osuModels.ScoreResponse, analyzeChan chan<- models.Score, db *gorm.DB, osuSvc *osuservices.OsuService, redisClient *database.RedisClient) {
	playerRepo := repositories.NewPlayerRepository(db)

	batch := make([]models.Score, 0, 50)
	queuedPlayers := make(map[uint]struct{})

	go func() {
		for {
			select {
			case scoreToFilter := <-filterChan:
				mappedScore, err := models.MapOsuScoreToModel(scoreToFilter)
				if err != nil {
					continue
				}

				if player, ok, err := redisClient.GetPlayer(mappedScore.PlayerID); err != nil {
					continue
				} else if ok {
					if applyFilters(player, []PlayerFilter{isFrench, isTop10k}) {
						fmt.Println("Found score for player", player.ID)
						analyzeChan <- mappedScore
					}
					continue
				}

				if _, alreadyQueued := queuedPlayers[mappedScore.PlayerID]; alreadyQueued {
					continue
				}

				queuedPlayers[mappedScore.PlayerID] = struct{}{}
				batch = append(batch, mappedScore)

				if len(batch) == 50 {
					checkPlayersFromScores(batch, analyzeChan, playerRepo, osuSvc, redisClient)
					for _, score := range batch {
						delete(queuedPlayers, score.PlayerID)
					}
					time.Sleep(3000 * time.Millisecond)
					batch = batch[:0]
				}

			case <-stopChan:
				return
			}
		}
	}()
}

func checkPlayersFromScores(scores []models.Score, analyzeChan chan<- models.Score, playerRepo *repositories.PlayerRepository, osuSvc *osuservices.OsuService, redisClient *database.RedisClient) {
	playerIds := make([]uint, 0, len(scores))
	for _, score := range scores {
		playerIds = append(playerIds, score.PlayerID)
	}

	request := osuModels.UsersRequest{IDs: playerIds}
	players, err := osuSvc.GetPlayers(request)
	if err != nil {
		return
	}

	if len(players.Users) == 0 {
		return
	}

	for _, player := range players.Users {
		newPlayer := models.MapOsuUserToModel(player)
		if applyFilters(newPlayer, []PlayerFilter{isFrench, isTop10k}) {
			err := playerRepo.Update(&newPlayer)
			if err != nil {
				return
			}

			for _, score := range scores {
				if score.PlayerID == player.ID {
					fmt.Println("Found score for player", player.ID)
					analyzeChan <- score
				}
			}
		}

		if err := redisClient.CachePlayer(newPlayer, time.Hour); err != nil {
			fmt.Println("Failed to cache player:", err)
		}
	}
}

type PlayerFilter func(player models.Player) bool

func isFrench(player models.Player) bool {
	return player.Country == "FR"
}

func isTop10k(player models.Player) bool {
	return player.GlobalRank <= 10000 && player.GlobalRank > 0
}

func applyFilters(player models.Player, filters []PlayerFilter) bool {
	for _, filter := range filters {
		if !filter(player) {
			return false
		}
	}
	return true
}
