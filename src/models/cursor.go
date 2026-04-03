package models

import "gorm.io/gorm"

type Cursor struct {
	gorm.Model
	Cursor string
}
