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

func (r *PlayerRepository) GetPlayers(page int) (models.PlayersResponse, error) {
	var players []models.Player
	if err := r.db.Order("global_rank ASC").Offset((page - 1) * 50).Limit(50).Find(&players).Error; err != nil {
		return models.PlayersResponse{}, err
	}

	var total int64
	if err := r.db.Model(&models.Player{}).Count(&total).Error; err != nil {
		return models.PlayersResponse{}, err
	}

	return models.PlayersResponse{
		Players: players,
		Page:    page,
		Total:   total,
	}, nil
}
