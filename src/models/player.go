package models

type Player struct {
	PlayerId   uint `gorm:"primaryKey"`
	Username   string
	GlobalRank uint
	Pp         float64
}

func MapOsuUserToModel(u User) Player {
	return Player{
		PlayerId:   u.ID,
		Username:   u.Username,
		GlobalRank: u.Statistics.GlobalRank,
		Pp:         u.Statistics.Pp,
	}
}
