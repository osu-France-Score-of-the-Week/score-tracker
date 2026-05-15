package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"score-tracker/models"
	"score-tracker/models/requests"
	"score-tracker/osuservices"
	"score-tracker/repositories"
	"time"

	"gorm.io/gorm"
)

func CreateAnalyzeScores(scoresChan chan<- models.Score, analyzeChan <-chan models.Score, stopChanAnalyzeScores <-chan struct{}, db *gorm.DB, osuSvc *osuservices.OsuService) {
	beatmapRepo := repositories.NewBeatmapRepository(db)
	beatmapsetRepo := repositories.NewBeatmapsetRepository(db)
	beatmapAttributes := repositories.NewBeatmapAttributesRepository(db)

	go func() {
		for {
			select {
			case score := <-analyzeChan:
				log.Printf("[analyze] score=%d player=%d beatmap=%d acc=%.6f pp=%.2f mods=%v stats=%v", score.ID, score.PlayerID, score.BeatmapID, score.Accuracy, score.Pp, score.Mods, score.Statistics)
				beatmap, attributes, err := UpdateBeatmap(score.BeatmapID, score.Mods, beatmapRepo, beatmapsetRepo, beatmapAttributes, osuSvc)
				if err != nil {
					log.Printf("[analyze] update beatmap failed beatmap=%d err=%v", score.BeatmapID, err)
					continue
				}
				log.Printf("[analyze] beatmap=%d cs=%.2f od=%.2f ar=%.2f version=%q attrs_star=%.3f attrs_max_combo=%d", beatmap.ID, beatmap.CS, beatmap.OD, beatmap.AR, beatmap.Version, attributes.StarRating, attributes.MaxCombo)

				result, err := RequestAnalyzeService(score, beatmap, attributes)
				if err != nil {
					log.Printf("[analyze] request failed score=%d beatmap=%d err=%v", score.ID, score.BeatmapID, err)
					continue
				}

				score.AnalyzedScore = result
				log.Printf("[analyze] score=%d result=%.6f", score.ID, result)

				scoresChan <- score
			case <-stopChanAnalyzeScores:
				return
			}
		}
	}()
}

func RequestAnalyzeService(
	score models.Score,
	beatmap models.Beatmap,
	attributes models.BeatmapAttributes,
) (float64, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	body, err := json.Marshal(requests.AnalyzeRequest{
		Beatmap:           beatmap,
		BeatmapAttributes: attributes,
		Score:             score,
	})
	if err != nil {
		log.Printf("[analyze] marshal failed score=%d beatmap=%d err=%v", score.ID, beatmap.ID, err)
		return 0, err
	}
	log.Printf("[analyze] sending request score=%d beatmap=%d body_bytes=%d", score.ID, beatmap.ID, len(body))

	req, err := http.NewRequest(
		"POST",
		"http://host.docker.internal:8180/analyze",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[analyze] http request failed score=%d beatmap=%d err=%v", score.ID, beatmap.ID, err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		log.Printf("[analyze] non-200 score=%d beatmap=%d status=%d body=%s", score.ID, beatmap.ID, resp.StatusCode, buf.String())
		return 0, fmt.Errorf("analyze service returned status %d", resp.StatusCode)
	}

	var result float64

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[analyze] decode failed score=%d beatmap=%d err=%v", score.ID, beatmap.ID, err)
		return 0, err
	}
	log.Printf("[analyze] response score=%d beatmap=%d result=%.6f", score.ID, beatmap.ID, result)

	return result, nil
}

func UpdateBeatmap(beatmapID uint, mods models.Mods, beatmapRepo *repositories.BeatmapRepository, beatmapsetRepo *repositories.BeatmapsetRepository, beatmapAttributes *repositories.BeatmapAttributesRepository, osuSvc *osuservices.OsuService) (models.Beatmap, models.BeatmapAttributes, error) {
	osuBeatmap, err := osuSvc.GetBeatmap(beatmapID)
	if err != nil {
		return models.Beatmap{}, models.BeatmapAttributes{}, err
	}

	beatmap := models.MapOsuBeatmapToModel(osuBeatmap)

	if err := beatmapsetRepo.Update(&beatmap.Beatmapset); err != nil {
		return models.Beatmap{}, models.BeatmapAttributes{}, err
	}

	if err := beatmapRepo.Update(&beatmap); err != nil {
		return models.Beatmap{}, models.BeatmapAttributes{}, err
	}

	osuAttributes, err := osuSvc.GetBeatmapAttribute(beatmapID, mods)
	if err != nil {
		return models.Beatmap{}, models.BeatmapAttributes{}, err
	}

	attributes := models.BeatmapAttributes{
		BeatmapID:                 beatmapID,
		Mods:                      mods,
		StarRating:                osuAttributes.Attributes.StarRating,
		MaxCombo:                  osuAttributes.Attributes.MaxCombo,
		AimDifficulty:             osuAttributes.Attributes.AimDifficulty,
		AimDifficultSliderCount:   osuAttributes.Attributes.AimDifficultSliderCount,
		SpeedDifficulty:           osuAttributes.Attributes.SpeedDifficulty,
		SpeedNoteCount:            osuAttributes.Attributes.SpeedNoteCount,
		SliderFactor:              osuAttributes.Attributes.SliderFactor,
		AimDifficultStrainCount:   osuAttributes.Attributes.AimDifficultStrainCount,
		SpeedDifficultStrainCount: osuAttributes.Attributes.SpeedDifficultStrainCount,
	}

	_, err = beatmapAttributes.GetByBeatmapModsCombination(beatmapID, mods)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return beatmap, attributes, beatmapAttributes.Create(&attributes)
		}
		return models.Beatmap{}, models.BeatmapAttributes{}, err
	}

	return beatmap, attributes, nil
}
