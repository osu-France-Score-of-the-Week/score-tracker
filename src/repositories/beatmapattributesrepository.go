package repositories

import (
	"score-tracker/models"

	"gorm.io/gorm"
)

type BeatmapAttributesRepository struct {
	db *gorm.DB
}

func NewBeatmapAttributesRepository(db *gorm.DB) *BeatmapAttributesRepository {
	return &BeatmapAttributesRepository{db: db}
}

func (r *BeatmapAttributesRepository) Create(beatmapAttributes *models.BeatmapAttributes) error {
	if err := r.db.Create(beatmapAttributes).Error; err != nil {
		return err
	}
	return nil
}

func (r *BeatmapAttributesRepository) Update(beatmapAttributes *models.BeatmapAttributes) error {
	if err := r.db.Save(beatmapAttributes).Error; err != nil {
		return err
	}
	return nil
}
