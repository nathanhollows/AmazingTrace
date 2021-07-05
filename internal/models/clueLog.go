package models

import (
	"time"

	"gorm.io/gorm"
)

// ClueLog keeps track of the current and found clues for each team
type ClueLog struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Team      string `gorm:"primaryKey"`
	ClueCode  string `gorm:"primaryKey"`
	Found     time.Time
	Active    bool
	Clue      Clue `gorm:"foreignKey:Code;references:ClueCode"`
}
