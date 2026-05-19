package jobs

import (
	"score-tracker/database"
	"score-tracker/models"
	osuModels "score-tracker/models/osu"
	"score-tracker/osuservices"
	"time"

	"gorm.io/gorm"
)

// TODO: Refactor everything, change structure
func CreateJobs(db *gorm.DB, osuSvc *osuservices.OsuService, redisClient *database.RedisClient) {
	stopChanRetrieveScores := make(chan struct{})
	stopChanFilter := make(chan struct{})
	stopChanAnalyzeScores := make(chan struct{})
	stopChanCreateScores := make(chan struct{})
	stopChanRecomputeScores := make(chan struct{})
	filterChan := make(chan osuModels.ScoreResponse, 2000)
	analyzeChan := make(chan models.Score, 100)
	scoresChan := make(chan models.Score, 100)

	RetrieveScores(20*time.Second, stopChanRetrieveScores, filterChan, osuSvc)
	FilterScores(stopChanFilter, filterChan, analyzeChan, db, osuSvc, redisClient)
	CreateAnalyzeScores(scoresChan, analyzeChan, stopChanAnalyzeScores, db, osuSvc)
	CreateScores(scoresChan, stopChanCreateScores, db, osuSvc)
	RecomputeScores(10*time.Second, stopChanRecomputeScores, db, osuSvc)
}
