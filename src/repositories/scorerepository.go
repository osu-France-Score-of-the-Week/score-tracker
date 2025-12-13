package repositories

import (
	"score-tracker/models"

	"gorm.io/gorm"
)

type ScoreRepository struct {
	db *gorm.DB
}

func NewScoreRepository(db *gorm.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

func (r *ScoreRepository) Create(score *models.Score) error {
	if err := r.db.Create(score).Error; err != nil {
		return err
	}
	return nil
}
