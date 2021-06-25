package models

import (
	"gorm.io/gorm"
)

// Clue stores a simple riddle based clue for a location
type Clue struct {
	gorm.Model
	Code      string `gorm:"index:unique"`
	Title     string
	Text      string
	Challenge bool
}
