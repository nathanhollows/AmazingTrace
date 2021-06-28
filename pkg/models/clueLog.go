package models

import "database/sql"

// ClueLog keeps track of the current and found clues for each team
type ClueLog struct {
	Team     string       `gorm:"primaryKey"`
	ClueCode string       `gorm:"primaryKey"`
	Found    sql.NullTime `gorm:"default:null"`
	Order    int          `gorm:"sort:asc"`
	Clue     Clue         `gorm:"foreignKey:Code;references:ClueCode"`
}
