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
	stopChanFilter := make(chan struct{})
	stopChanAnalyzeScores := make(chan struct{})
	stopChanCreateScores := make(chan struct{})
	stopChanRecomputeScores := make(chan struct{})
	filterChan := make(chan osuModels.ScoreResponse, 2000)
	analyzeChan := make(chan models.Score, 100)
	scoresChan := make(chan models.Score, 100)

	RetrieveScores(20*time.Second, stopChanRetrieveScores, filterChan, osuSvc)
	FilterScores(stopChanFilter, filterChan, analyzeChan, db, osuSvc)
	CreateAnalyzeScores(scoresChan, analyzeChan, stopChanAnalyzeScores, db, osuSvc)
	CreateScores(scoresChan, stopChanCreateScores, db, osuSvc)
	RecomputeScores(2*time.Second, stopChanRecomputeScores, db, osuSvc)
}
