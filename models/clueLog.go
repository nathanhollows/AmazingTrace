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

type Clues []ClueLog

// Return the set of active clues for a team
func (c Clues) Active() Clues {
	var active Clues
	for _, clue := range c {
		if clue.Active {
			active = append(active, clue)
		}
	}
	return active
}

// Return the set of found clues for a team
func (c Clues) Found() Clues {
	var found Clues
	for _, clue := range c {
		if !clue.Found.IsZero() {
			found = append(found, clue)
		}
	}
	return found
}
