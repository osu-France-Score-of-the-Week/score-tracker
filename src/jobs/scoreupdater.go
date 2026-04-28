package jobs

import (
	"fmt"
	osuModels "score-tracker/models/osu"
	"score-tracker/osuservices"
	"time"
)

func RetrieveScores(interval time.Duration, stopChan <-chan struct{}, filterChan chan<- osuModels.ScoreResponse, osuSvc *osuservices.OsuService) {

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		cursor := ""
		for {
			select {
			case <-ticker.C:
				recentScores, err := osuSvc.GetRecentScores(cursor)
				if err != nil {
					fmt.Println("Error updating recent scores:", err)
					continue
				}

				for _, score := range recentScores.Scores {
					filterChan <- score
				}

				cursor = recentScores.Cursor

				fmt.Println("new cursor:", recentScores.Cursor)

			case <-stopChan:
				return
			}
		}
	}()
}
