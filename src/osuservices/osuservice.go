package osuservices

import (
	"score-tracker/queue"
	"time"

	"gorm.io/gorm"
)

type OsuService struct {
	db    *gorm.DB
	queue *queue.RequestQueue
}

func NewOsuService(db *gorm.DB) *OsuService {
	return &OsuService{
		db:    db,
		queue: queue.NewQueue(2 * time.Second),
	}
}
