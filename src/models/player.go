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

type PlayerResponse struct {
	ID         uint
	Username   string
	GlobalRank uint
	Pp         float64
	ScoreCount uint `gorm:"-" json:"scoreCount"`
}

type PlayersPageResponse struct {
	Players    []PlayerResponse `json:"players"`
	Page       int              `json:"page"`
	TotalPages int64            `json:"totalPages"`
}
