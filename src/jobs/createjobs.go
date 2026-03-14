package jobs

import (
	"score-tracker/models"
	"time"

	"gorm.io/gorm"
)

func CreateJobs(db *gorm.DB) {
	stopChanRetrieveScores := make(chan struct{})
	stopChanFilter := make(chan struct{}, 1000)
	stopChanCreateScores := make(chan struct{}, 1000)
	filterChan := make(chan models.OsuScore, 1000)
	scoresChan := make(chan models.Score, 100)
	RetrieveScores(5*time.Second, stopChanRetrieveScores, filterChan)
	FilterScores(stopChanFilter, filterChan, scoresChan)
	CreateScores(scoresChan, stopChanCreateScores, db)
}
