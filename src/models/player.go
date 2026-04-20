package models

type Player struct {
	ID         uint `gorm:"primaryKey"`
	Username   string
	GlobalRank uint
	Pp         float64
}

func MapOsuUserToModel(u User) Player {
	return Player{
		ID:         u.ID,
		Username:   u.Username,
		GlobalRank: u.Statistics.GlobalRank,
		Pp:         u.Statistics.Pp,
	}
}

type PlayersResponse struct {
	Players []Player `json:"players"`
	Page    int      `json:"page"`
	Total   int64    `json:"total"`
}
