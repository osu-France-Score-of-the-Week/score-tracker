package jobs

import (
	"score-tracker/models"
	"score-tracker/osuservices"
	"time"

	"gorm.io/gorm"
)

func CreateJobs(db *gorm.DB, osuSvc *osuservices.OsuService) {
	stopChanRetrieveScores := make(chan struct{})
	stopChanFilter := make(chan struct{}, 2000)
	stopChanCreateScores := make(chan struct{}, 2000)
	filterChan := make(chan models.OsuScore, 2000)
	scoresChan := make(chan models.Score, 100)

	RetrieveScores(20*time.Second, stopChanRetrieveScores, filterChan, osuSvc)
	FilterScores(stopChanFilter, filterChan, scoresChan, db, osuSvc)
	CreateScores(scoresChan, stopChanCreateScores, db, osuSvc)
}
