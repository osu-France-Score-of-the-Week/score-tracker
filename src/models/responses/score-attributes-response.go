package responses

import (
	"score-tracker/models"
)

type ScoreWithAttributes struct {
	ID         uint
	Accuracy   float64
	BeatmapID  uint
	Beatmap    models.Beatmap
	EndedAt    int64
	HasReplay  bool
	MaxCombo   uint
	Mods       models.Mods `gorm:"type:jsonb"`
	Pp         float64
	Rank       string
	Statistics models.ScoreStatistics `gorm:"type:jsonb"`
	PlayerID   uint
	Player     Player
	Attributes models.BeatmapAttributes `gorm:"-"`
}
