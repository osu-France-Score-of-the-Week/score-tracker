package repositories

import (
	"score-tracker/models"

	"gorm.io/gorm"
)

type CursorRepository struct {
	db *gorm.DB
}

func NewCursorRepository(db *gorm.DB) *CursorRepository {
	return &CursorRepository{db: db}
}

func (r *CursorRepository) Create(cursor *models.Cursor) error {
	if err := r.db.Create(cursor).Error; err != nil {
		return err
	}
	return nil
}

func (r *CursorRepository) GetLastCursor() (models.Cursor, error) {
	var cursor models.Cursor
	if err := r.db.Order("created_at desc").First(&cursor).Error; err != nil {
		return models.Cursor{}, err
	}
	return cursor, nil
}
