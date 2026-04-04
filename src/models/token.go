package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	Token     string
	ExpiresAt time.Time
}
