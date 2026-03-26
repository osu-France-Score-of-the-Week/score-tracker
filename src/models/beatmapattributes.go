package models

type BeatmapAttributes struct {
	ID         uint `gorm:"primaryKey;column:id"`
	BeatmapID  uint
	Beatmap    Beatmap
	Mods       Mods `gorm:"type:jsonb"`
	StarRating float64
	MaxCombo   uint
}
