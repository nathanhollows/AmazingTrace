package models

import "gorm.io/gorm"

// Code is a list of Codes for the scavenger hunt mode
type Code struct {
	gorm.Model
	Code string
}
