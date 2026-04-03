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

func (r *BeatmapAttributesRepository) GetByBeatmapModsCombination(beatmapID uint, mods models.Mods) (*models.BeatmapAttributes, error) {
	var beatmapAttributes models.BeatmapAttributes
	if err := r.db.Where("beatmap_id = ? AND mods = ?", beatmapID, mods).First(&beatmapAttributes).Error; err != nil {
		return nil, err
	}
	return &beatmapAttributes, nil
}
