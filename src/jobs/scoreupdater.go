package jobs

import (
	"fmt"
	"score-tracker/models"
	"score-tracker/osuservices"
	"score-tracker/repositories"
	"time"

	"gorm.io/gorm"
)

func RetrieveScores(interval time.Duration, stopChan <-chan struct{}, filterChan chan<- models.OsuScore, db *gorm.DB) {
	cursorRepo := repositories.NewCursorRepository(db)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cursordb, err := cursorRepo.GetLastCursor()
				// fmt.Println("Updating recent scores...")
				recentScores, err := osuservices.GetRecentScores(cursordb.Cursor)
				if err != nil {
					fmt.Println("Error updating recent scores:", err)
					continue
				}

				for _, score := range recentScores.Scores {
					filterChan <- score
				}

				err = cursorRepo.Create(&models.Cursor{
					Cursor: recentScores.Cursor,
				})

				if err != nil {
					fmt.Println("Error saving cursor to database:", err)
					return
				}

			case <-stopChan:
				return
			}
		}
	}()
}
