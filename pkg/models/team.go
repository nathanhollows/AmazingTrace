package models

import (
	"github.com/nathanhollows/AmazingTrace/pkg/helpers"
	"gorm.io/gorm"
)

// Team holds the team specific information for a given team
type Team struct {
	gorm.Model
	Code     string `gorm:"uniqueIndex:idx_code,sort:desc;not null"`
	Assigned bool   `gorm:"default:false"`
	Started  bool   `gorm:"default:false;not null"`
}

// BeforeCreate generates a random string for the team to identify by
func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	t.Code = helpers.NewCode(4)
	return
}
