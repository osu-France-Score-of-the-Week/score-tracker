package jobs

import (
	"fmt"
	"score-tracker/models"
	"score-tracker/repositories"

	"gorm.io/gorm"
)

func CreateScores(scoresChan <-chan models.OsuScore, stopChan <-chan struct{}, db *gorm.DB) {
	scoreRepo := repositories.NewScoreRepository(db)

	go func() {
		for {
			select {
			case score := <-scoresChan:
				scoreCreate, err := models.MapOsuScoreToModel(score)
				if err != nil {
					fmt.Println("Error mapping OsuScore to Score model:", err)
					continue
				}

				if err := scoreRepo.Create(&scoreCreate); err != nil {
					fmt.Println("Error creating score in database:", err)
					continue
				}

			case <-stopChan:
				return
			}
		}
	}()
}
