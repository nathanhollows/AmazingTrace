package models

import (
	"database/sql"
	"time"

	"github.com/nathanhollows/AmazingTrace/pkg/helpers"
	"gorm.io/gorm"
)

// Clue stores a simple riddle based clue for a location
type Clue struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	Code      string `gorm:"uniqueIndex:clue_code,sort:desc;not null;primarykey"`
	Location  string `gorm:"index:unique;not null"`
	Clue      string `gorm:"not null"`
}

// BeforeCreate generates a random string for the clue to be identified by
func (c *Clue) BeforeCreate(tx *gorm.DB) (err error) {
	c.Code = helpers.NewCode(4)
	return
}
