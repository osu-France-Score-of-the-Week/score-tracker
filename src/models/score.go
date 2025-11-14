package models

type Score struct {
	ScoreId    uint    `gorm:"primaryKey"`
	Accuracy   float64 `gorm:"not null"`
	BeatmapId  uint    `gorm:"not null"`
	EndedAt    int64   `gorm:"not null"`
	HasReplay  bool    `gorm:"not null"`
	MaxCombo   uint    `gorm:"not null"`
	Mods       string  `gorm:"not null"`
	Pp         float64 `gorm:"not null"`
	Rank       uint    `gorm:"not null"`
	Statistics string  `gorm:"not null"`
	PlayerId   uint    `gorm:"not null"`
}
