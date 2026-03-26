package repositories

import (
	"score-tracker/models"

	"gorm.io/gorm"
)

type BeatmapRepository struct {
	db *gorm.DB
}

func NewBeatmapRepository(db *gorm.DB) *BeatmapRepository {
	return &BeatmapRepository{db: db}
}

func (r *BeatmapRepository) Create(beatmap *models.Beatmap) error {
	if err := r.db.Create(beatmap).Error; err != nil {
		return err
	}
	return nil
}

func (r *BeatmapRepository) Update(beatmap *models.Beatmap) error {
	if err := r.db.Save(beatmap).Error; err != nil {
		return err
	}
	return nil
}
