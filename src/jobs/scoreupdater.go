package jobs

import (
	"fmt"
	"score-tracker/models"
	"score-tracker/osuservices"
	"time"
)

func RetrieveScores(interval time.Duration, stopChan <-chan struct{}, filterChan chan<- models.OsuScore) {

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		cursor := ""
		for {
			select {
			case <-ticker.C:
				recentScores, err := osuservices.GetRecentScores(cursor)
				if err != nil {
					fmt.Println("Error updating recent scores:", err)
					continue
				}

				for _, score := range recentScores.Scores {
					filterChan <- score
				}

				cursor = recentScores.Cursor

				if err != nil {
					fmt.Println("Error saving cursor to database:", err)
					return
				}

				fmt.Println("new cursor:", recentScores.Cursor)

			case <-stopChan:
				return
			}
		}
	}()
}
