package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Admin stores logins for the admin panel
type Admin struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateAdmin creates an admin user
func CreateAdmin(db gorm.DB, username string, password string) (Admin, error) {
	if username == "" || password == "" {
		return Admin{}, errors.New("ADMIN_USERNAME and ADMIN_PASSWORD should be set")
	}
	if username == "admin" || password == "password" {
		return Admin{}, errors.New("ADMIN_USERNAME and ADMIN_PASSWORD should be changed")
	}
	var admin Admin
	db.Model(&admin).Where("username = ?", username).Find(&admin)
	if admin.Username == "" {
		admin.Username = username
		admin.Password, _ = admin.HashPassword(password)
		db.Create(&admin)
	}

	return admin, nil
}

// CheckHashPassword checks a password against a hash
func (a Admin) CheckHashPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// HashPassword hashes a password
func (a Admin) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
