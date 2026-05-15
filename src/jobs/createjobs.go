package jobs

import (
	"score-tracker/models"
	osuModels "score-tracker/models/osu"
	"score-tracker/osuservices"
	"time"

	"gorm.io/gorm"
)

func CreateJobs(db *gorm.DB, osuSvc *osuservices.OsuService) {
	stopChanRetrieveScores := make(chan struct{})
	stopChanFilter := make(chan struct{}, 2000)
	stopChanAnalyzeScores := make(chan struct{}, 2000)
	stopChanCreateScores := make(chan struct{}, 2000)
	filterChan := make(chan osuModels.ScoreResponse, 2000)
	analyzeChan := make(chan models.Score, 100)
	scoresChan := make(chan models.Score, 100)

	RetrieveScores(20*time.Second, stopChanRetrieveScores, filterChan, osuSvc)
	FilterScores(stopChanFilter, filterChan, analyzeChan, db, osuSvc)
	CreateAnalyzeScores(scoresChan, analyzeChan, stopChanAnalyzeScores, db, osuSvc)
	CreateScores(scoresChan, stopChanCreateScores, db, osuSvc)
}
