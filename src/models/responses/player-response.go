package responses

type Player struct {
	ID         uint
	Username   string
	GlobalRank uint
	Pp         float64
	ScoreCount uint `gorm:"-" json:"scoreCount"`
}
