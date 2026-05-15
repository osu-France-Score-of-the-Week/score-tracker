package jobs

import (
	"errors"
	"log"
	"score-tracker/osuservices"
	"score-tracker/repositories"
	"time"

	"gorm.io/gorm"
)

func RecomputeScores(interval time.Duration, stopChan <-chan struct{}, db *gorm.DB, osuSvc *osuservices.OsuService) {
	scoresRepo := repositories.NewScoreRepository(db)
	beatmapRepo := repositories.NewBeatmapRepository(db)
	beatmapsetRepo := repositories.NewBeatmapsetRepository(db)
	beatmapAttributesRepo := repositories.NewBeatmapAttributesRepository(db)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				score, err := scoresRepo.GetOneScoreWithMissingAnalyses()
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						continue
					}
					log.Println(err)
					continue
				}

				if score != nil {
					log.Println("Recomputing score:", score)
					beatmap, attributes, err := UpdateBeatmap(score.BeatmapID, score.Mods, beatmapRepo, beatmapsetRepo, beatmapAttributesRepo, osuSvc)
					if err != nil {
						log.Println(err)
						continue
					}
					analyzedScore, err := RequestAnalyzeService(*score, beatmap, attributes)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Println("Analyzed score:", analyzedScore)
					score.AnalyzedScore = analyzedScore
					if err := scoresRepo.UpdateScore(*score); err != nil {
						log.Println(err)
					}
				}
			case <-stopChan:
				return
			}
		}
	}()
}
