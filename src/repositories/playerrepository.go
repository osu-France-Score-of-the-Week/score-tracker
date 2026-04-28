package repositories

import (
	"score-tracker/models"
	"score-tracker/models/responses"

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

func (r *PlayerRepository) GetPlayers(page int) (responses.PlayersPage, error) {
	const pageSize int64 = 50

	var players []responses.Player
	if err := r.db.Model(&models.Player{}).Order("global_rank ASC").Offset((page - 1) * int(pageSize)).Limit(int(pageSize)).Find(&players).Error; err != nil {
		return responses.PlayersPage{}, err
	}

	var totalPlayers int64
	if err := r.db.Model(&models.Player{}).Count(&totalPlayers).Error; err != nil {
		return responses.PlayersPage{}, err
	}

	totalPages := (totalPlayers + pageSize - 1) / pageSize

	for i := range players {
		var scoreCount int64
		if err := r.db.Model(&models.Score{}).Where("player_id = ?", players[i].ID).Count(&scoreCount).Error; err != nil {
			return responses.PlayersPage{}, err
		}
		players[i].ScoreCount = uint(scoreCount)
	}

	return responses.PlayersPage{
		Players:    players,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}
