package jobs

import (
	"fmt"
	"score-tracker/models"
	"score-tracker/osuservices"
	"time"
)

func RetrieveScores(interval time.Duration, stopChan <-chan struct{}, scoresChan chan<- models.OsuScore) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		var cursor *string

		for {
			select {
			case <-ticker.C:
				fmt.Println("Updating recent scores...")
				recentScores, err := osuservices.GetRecentScores(cursor)
				if err != nil {
					fmt.Println("Error updating recent scores:", err)
					continue
				}

				go func() {
					for _, score := range recentScores.Scores {
						scoresChan <- score
					}
				}()

				cursor = &recentScores.Cursor

			case <-stopChan:
				return
			}
		}
	}()
}
