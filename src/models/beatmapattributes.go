package models

type BeatmapAttributes struct {
	ID         uint `gorm:"primaryKey;column:id"`
	BeatmapID  uint `gorm:"uniqueIndex:idx_beatmap_mods;column:beatmap_id"`
	Beatmap    Beatmap
	Mods       Mods `gorm:"type:jsonb;uniqueIndex:idx_beatmap_mods"`
	StarRating float64
	MaxCombo   uint
}
