package models

import "gorm.io/gorm"

// Area is a physical location for clues
type Area struct {
	gorm.Model
	Name string
}
