package osuservices

import "gorm.io/gorm"

type OsuService struct {
	db *gorm.DB
}

func NewOsuService(db *gorm.DB) *OsuService {
	return &OsuService{db: db}
}
