package repositories

import (
	"score-tracker/models"

	"gorm.io/gorm"
)

type BeatmapsetRepository struct {
	db *gorm.DB
}

func NewBeatmapsetRepository(db *gorm.DB) *BeatmapsetRepository {
	return &BeatmapsetRepository{db: db}
}

func (r *BeatmapsetRepository) Create(beatmapset *models.Beatmapset) error {
	if err := r.db.Create(beatmapset).Error; err != nil {
		return err
	}
	return nil
}

func (r *BeatmapsetRepository) Update(beatmapset *models.Beatmapset) error {
	if err := r.db.Save(beatmapset).Error; err != nil {
		return err
	}
	return nil
}
