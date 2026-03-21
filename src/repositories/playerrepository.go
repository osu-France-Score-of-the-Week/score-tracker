package repositories

import (
	"score-tracker/models"

	"gorm.io/gorm"
)

type PlayerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Create(player *models.Player) error {
	if err := r.db.Create(player).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlayerRepository) Update(player *models.Player) error {
	if err := r.db.Save(player).Error; err != nil {
		return err
	}
	return nil
}
