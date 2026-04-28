package models

import "score-tracker/models/osu"

type Player struct {
	ID         uint `gorm:"primaryKey"`
	Username   string
	GlobalRank uint
	Pp         float64
}

func MapOsuUserToModel(u osu.UserResponse) Player {
	return Player{
		ID:         u.ID,
		Username:   u.Username,
		GlobalRank: u.Statistics.GlobalRank,
		Pp:         u.Statistics.Pp,
	}
}
