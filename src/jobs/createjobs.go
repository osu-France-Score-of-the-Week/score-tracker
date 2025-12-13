package jobs

import (
	"score-tracker/models"
	"time"

	"gorm.io/gorm"
)

func CreateJobs(db *gorm.DB) {
	stopChanRetrieveScores := make(chan struct{})
	stopChanCreateScores := make(chan struct{}, 1000)
	scoresChan := make(chan models.OsuScore, 100)
	RetrieveScores(5*time.Second, stopChanRetrieveScores, scoresChan)
	CreateScores(scoresChan, stopChanCreateScores, db)
}
