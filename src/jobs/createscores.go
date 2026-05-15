package jobs

import (
	"fmt"
	"score-tracker/models"
	"score-tracker/osuservices"
	"score-tracker/repositories"

	"gorm.io/gorm"
)

func CreateScores(scoresChan <-chan models.Score, stopChan <-chan struct{}, db *gorm.DB, osuSvc *osuservices.OsuService) {
	scoreRepo := repositories.NewScoreRepository(db)

	go func() {
		for {
			select {
			case score := <-scoresChan:
				if err := scoreRepo.Create(&score); err != nil {
					fmt.Println("Error creating score in database:", err)
					continue
				}
				fmt.Println("Created score in database:", score.ID)

			case <-stopChan:
				return
			}
		}
	}()
}
