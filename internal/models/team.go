package models

import (
	"errors"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/nathanhollows/AmazingTrace/internal/helpers"
	"gorm.io/gorm"
)

// Team holds the team specific information for a given team
type Team struct {
	gorm.Model
	Code     string    `gorm:"uniqueIndex:idx_code,sort:desc;not null"`
	Assigned bool      `gorm:"default:false"`
	Started  bool      `gorm:"default:false;not null"`
	ClueLog  []ClueLog `gorm:"foreignKey:Team;references:Code"`
	Found    int       `gorm:"default:0;not null"`
}

// BeforeCreate generates a random string for the team to identify by
func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	t.Code = helpers.NewCode(4)
	return
}

// AfterFind generates a random string for the team to identify by
func (t *Team) AfterFind(tx *gorm.DB) (err error) {
	sort.SliceStable(t.ClueLog, func(i, j int) bool {
		if t.ClueLog[i].Active != t.ClueLog[j].Active { // If one of the clues is active
			return t.ClueLog[i].Active
		}
		if t.ClueLog[i].Found.IsZero() {
			return t.ClueLog[i].Clue.Location < t.ClueLog[j].Clue.Location
		}
		return t.ClueLog[i].Found.Before(t.ClueLog[j].Found)
	})
	return
}

// Start creates a set of clues for the team
func (t *Team) Start(tx *gorm.DB) (err error) {
	t.Started = true

	clues := []Clue{}
	tx.Where("1=1").Find(&clues)

	for _, clue := range clues {
		clueLog := ClueLog{
			Team:     t.Code,
			ClueCode: clue.Code,
			Active:   false,
		}
		tx.Create(&clueLog)
	}

	t.ActivateClues(tx, 3)

	return
}

// ActivateClues ensures that x number of random clues are active, where possible
func (t *Team) ActivateClues(tx *gorm.DB, count int64) (err error) {
	t.Started = true

	// Check if we already have the right number of clues
	var activeCount int64
	tx.Where("team = ? AND active = true", t.Code).Find(&ClueLog{}).Count(&activeCount)
	if activeCount == count {
		return nil
	}

	// Otherwise, let's find the clues we can work with
	potentialClues := []ClueLog{}
	result := tx.Model(&ClueLog{}).
		Where("team = ? AND active = false AND found = ?", t.Code, time.Time{}).
		Find(&potentialClues)

	// If there are fewer clues available than the count then activate them all
	if (result.RowsAffected + activeCount) <= count {
		for _, clue := range potentialClues {
			tx.Model(&ClueLog{}).
				Where("clue_code = ? and team = ?", clue.ClueCode, clue.Team).
				Updates(map[string]interface{}{"active": true}).
				Omit("Clue")
		}
		return nil
	}

	// Otherwise create a random set of clues
	picks := make(map[int]bool)
	for len(picks) < int(count-activeCount) {
		pick := int(math.Floor(rand.Float64() * float64(result.RowsAffected)))
		picks[pick] = true
	}

	for pick := range picks {
		tx.Model(&ClueLog{}).
			Where("clue_code = ? and team = ?", potentialClues[pick].ClueCode, t.Code).
			Updates(map[string]interface{}{"active": true}).
			Omit("Clue")
	}

	return nil
}

// CheckClue checks if the clue is valid, and solvable
func (t *Team) CheckClue(tx *gorm.DB, code string) (err error) {
	// Check if the clue is valid
	code = strings.ToUpper(code)
	clue := &Clue{Code: code}
	result := tx.Find(&clue)
	if result.RowsAffected == 0 {
		return errors.New("That clue isn't valid")
	}

	// Check if the clue is valid, and solvable
	for _, log := range t.ClueLog {
		// If the clue doesn't match then skip it
		if log.Clue.Code != clue.Code {
			continue
		}
		// If the clue has already been found then return an error
		if !log.Found.IsZero() {
			return errors.New("You have already solved this clue")
		}
		// If the clue has just been solved then update
		if log.Active {
			tx.Model(&ClueLog{}).
				Where("clue_code = ? and team = ?", log.ClueCode, log.Team).
				Updates(map[string]interface{}{"found": time.Now(), "active": false}).
				Omit("Clue")
			tx.Model(&Team{}).
				Where("code = ?", t.Code).
				Update("found", t.Found+1)
			t.ActivateClues(tx, 3)
			return nil
		}
	}

	return errors.New("Nice try. That's not a clue you can solve yet")
}
