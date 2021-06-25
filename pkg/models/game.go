package models

import (
	"time"

	"gorm.io/gorm"
)

// Game stores the start and end times for each play through
type Game struct {
	gorm.Model
	StartTime time.Time
	EndTime   time.Time
}
