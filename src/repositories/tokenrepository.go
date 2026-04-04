package repositories

import (
	"score-tracker/models"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *TokenRepository) GetValidToken() (*models.Token, error) {
	var token models.Token
	if err := r.db.Where("expires_at > ?", gorm.Expr("NOW()")).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}
