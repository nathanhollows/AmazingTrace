package models

import "gorm.io/gorm"

// User represents both players and admins
type User struct {
	gorm.Model
	Email     string `json:"email"`
	Code      string `json:"code"`
	Admin     bool   `json:"admin"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	Temporary string `json:"temporary"`
}
