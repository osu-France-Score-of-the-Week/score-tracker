package jobs

import (
	"fmt"
	"score-tracker/models"
	"score-tracker/osuservices"
	"score-tracker/repositories"

	"gorm.io/gorm"
)

func CreateScores(scoresChan <-chan models.Score, stopChan <-chan struct{}, db *gorm.DB) {
	scoreRepo := repositories.NewScoreRepository(db)
	beatmapRepo := repositories.NewBeatmapRepository(db)
	beatmapsetRepo := repositories.NewBeatmapsetRepository(db)
	beatmapAttributes := repositories.NewBeatmapAttributesRepository(db)

	go func() {
		for {
			select {
			case score := <-scoresChan:
				err := UpdateBeatmap(score.BeatmapID, score.Mods, beatmapRepo, beatmapsetRepo, beatmapAttributes)
				if err != nil {
					fmt.Println("Error updating beatmap:", err)
					continue
				}

				if err := scoreRepo.Create(&score); err != nil {
					fmt.Println("Error creating score in database:", err)
					continue
				}
				fmt.Println("Created score in database:", score.ID)

			case <-stopChan:
				return
			}
		}
	}()
}

func UpdateBeatmap(beatmapID uint, mods models.Mods, beatmapRepo *repositories.BeatmapRepository, beatmapsetRepo *repositories.BeatmapsetRepository, beatmapAttributes *repositories.BeatmapAttributesRepository) error {
	osuBeatmap, err := osuservices.GetBeatmap(beatmapID)
	if err != nil {
		return err
	}

	beatmap := models.MapOsuBeatmapToModel(osuBeatmap)

	if err := beatmapsetRepo.Update(&beatmap.Beatmapset); err != nil {
		return err
	}

	if err := beatmapRepo.Update(&beatmap); err != nil {
		return err
	}

	osuAttributes, err := osuservices.GetBeatmapAttribute(beatmapID, mods)
	if err != nil {
		return err
	}

	attributes := models.BeatmapAttributes{
		BeatmapID:  beatmapID,
		Mods:       mods,
		StarRating: osuAttributes.Attributes.StarRating,
		MaxCombo:   osuAttributes.Attributes.MaxCombo,
	}

	_, err = beatmapAttributes.GetByBeatmapModsCombination(beatmapID, mods)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return beatmapAttributes.Create(&attributes)
		}
		return nil
	}

	return nil
}
