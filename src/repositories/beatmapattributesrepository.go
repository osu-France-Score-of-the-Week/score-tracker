package repositories

import (
	"score-tracker/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *BeatmapAttributesRepository) Upsert(beatmapAttributes *models.BeatmapAttributes) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "beatmap_id"},
			{Name: "mods"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"star_rating",
			"max_combo",
			"aim_difficulty",
			"aim_difficult_slider_count",
			"speed_difficulty",
			"speed_note_count",
			"slider_factor",
			"aim_difficult_strain_count",
			"speed_difficult_strain_count",
		}),
	}).Create(beatmapAttributes).Error
}

func (r *BeatmapAttributesRepository) GetByBeatmapModsCombination(beatmapID uint, mods models.Mods) (*models.BeatmapAttributes, error) {
	var beatmapAttributes models.BeatmapAttributes
	if err := r.db.Where("beatmap_id = ? AND mods = ?", beatmapID, mods).First(&beatmapAttributes).Error; err != nil {
		return nil, err
	}
	return &beatmapAttributes, nil
}
