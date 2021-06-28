package models

import (
	"database/sql"
	"sort"

	"github.com/nathanhollows/AmazingTrace/pkg/helpers"
	"gorm.io/gorm"
)

// Team holds the team specific information for a given team
type Team struct {
	gorm.Model
	Code     string    `gorm:"uniqueIndex:idx_code,sort:desc;not null"`
	Assigned bool      `gorm:"default:false"`
	Started  bool      `gorm:"default:false;not null"`
	ClueLog  []ClueLog `gorm:"foreignKey:Team;references:Code"`
}

// BeforeCreate generates a random string for the team to identify by
func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	t.Code = helpers.NewCode(4)
	return
}

// AfterFind generates a random string for the team to identify by
func (t *Team) AfterFind(tx *gorm.DB) (err error) {
	sort.SliceStable(t.ClueLog, func(i, j int) bool {
		return t.ClueLog[i].Order < t.ClueLog[j].Order
	})
	return
}

// Start creates a set of clues for the team
func (t *Team) Start(tx *gorm.DB) (err error) {
	t.Started = true

	clues := []Clue{}
	tx.Where("1=1").Find(&clues)

	for i, clue := range clues {
		clueLog := ClueLog{
			Team:     t.Code,
			ClueCode: clue.Code,
			Found:    sql.NullTime{},
			Order:    i,
		}
		tx.Create(&clueLog)
	}
	return
}
