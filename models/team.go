package models

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/nathanhollows/AmazingTrace/helpers"
	"gorm.io/gorm"
)

// Team holds the team specific information for a given team
type Team struct {
	gorm.Model
	Code     string `gorm:"uniqueIndex:idx_code,sort:desc;not null"`
	Name     string ``
	Assigned bool   `gorm:"default:false"`
	Started  bool   `gorm:"default:false;not null"`
	Clues    Clues  `gorm:"foreignKey:Team;references:Code"`
	Found    int    `gorm:"default:0;not null"`
	Delayed  bool   `gorm:"default:false"`
	Points   int    `gorm:"default:0;not null"`
}

// Get will fetch the team, and hydrate all associated records, given a team code.
// Returns both a team and an error.
func (t Team) Get(db *gorm.DB, code interface{}) (*Team, error) {
	team := Team{}
	code = strings.ToUpper(fmt.Sprintf("%v", code))
	result := db.Where("Code = ?", code).Preload("ClueLog.Clue").Find(&team)
	return &team, result.Error
}

// BeforeCreate generates a random string for the team to identify by
func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	t.Code = helpers.NewCode(4)
	return
}

// AfterFind generates a random string for the team to identify by
func (t *Team) AfterFind(tx *gorm.DB) (err error) {
	sort.SliceStable(t.Clues, func(i, j int) bool {
		if t.Clues[i].Active != t.Clues[j].Active { // If one of the clues is active
			return t.Clues[i].Active
		}
		if t.Clues[i].Found.IsZero() {
			return t.Clues[i].Clue.UpdatedAt.String() < t.Clues[j].Clue.UpdatedAt.String()
		}
		return t.Clues[i].Found.Before(t.Clues[j].Found)
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
			Found:    time.Time{},
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

// Solve will solve a clue for a team
func (t *Team) Solve(tx *gorm.DB, log ClueLog) error {
	tx.Model(&ClueLog{}).
		Where("clue_code = ? and team = ?", log.ClueCode, log.Team).
		Updates(map[string]interface{}{"found": time.Now(), "active": false}).
		Omit("Clue")
	t.UpdateCount(tx)
	t.ActivateClues(tx, 3)
	return nil
}

// UpdateCount will update the found number of the team
func (t *Team) UpdateCount(tx *gorm.DB) error {
	var count int64
	tx.Model(ClueLog{}).Where("team = ? AND found <> ?", t.Code, time.Time{}).Count(&count)
	tx.Model(&Team{}).
		Where("code = ?", t.Code).
		Update("found", count)
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
	for _, log := range t.Clues {
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
			t.Solve(tx, log)

			return nil
		}
	}

	return errors.New("Nice try. That's not a clue you can solve yet")
}

// FastForward automatically solves the next 3 clues
func (t *Team) FastForward(tx *gorm.DB, limit int) error {
	clues := []ClueLog{}
	// Solve the oldest clue first
	result := tx.Model(ClueLog{}).Where("team = ? AND active = 1", t.Code).Order("updated_at asc").Limit(limit).Find(&clues)
	if result.Error != nil {
		return errors.New("Could not find clues")
	}
	for _, clue := range clues {
		t.Solve(tx, clue)
	}
	return nil
}

// Rewind will find the last solved clue and unsolve it
func (t *Team) Rewind(tx *gorm.DB) error {
	lastClue := ClueLog{}
	result := tx.Model(ClueLog{}).Where("team = ? AND found <> ?", t.Code, time.Time{}).Order("found DESC").Limit(1).Find(&lastClue)
	if result.Error != nil {
		return errors.New("Could not find clue")
	} else if result.RowsAffected == 0 {
		return errors.New("There are no clues to rewind")
	}

	newestClue := ClueLog{}
	result = tx.Model(ClueLog{}).Where("team = ? AND active = 1", t.Code).Order("updated_at DESC").Limit(1).Find(&newestClue)
	if result.Error != nil {
		return errors.New("Could not find clue")
	}

	tx.Model(&ClueLog{}).
		Where("clue_code = ? and team = ?", newestClue.ClueCode, newestClue.Team).
		Updates(map[string]interface{}{"active": false}).
		Omit("Clue")
	tx.Model(&ClueLog{}).
		Where("clue_code = ? and team = ?", lastClue.ClueCode, lastClue.Team).
		Updates(map[string]interface{}{"active": true, "found": time.Time{}}).
		Omit("Clue")
	t.ActivateClues(tx, 3)
	t.UpdateCount(tx)

	return nil

}

// Shuffle will reset the clues a team is looking for
func (t *Team) Shuffle(tx *gorm.DB) error {
	clues := []ClueLog{}
	result := tx.Model(ClueLog{}).Where("team = ? AND active = 1", t.Code).Find(&clues)
	if result.Error != nil {
		return errors.New("Could not find clues")
	}
	for _, clue := range clues {
		tx.Model(&ClueLog{}).
			Where("clue_code = ? and team = ?", clue.ClueCode, clue.Team).
			Updates(map[string]interface{}{"active": false}).
			Omit("Clue")
	}
	t.ActivateClues(tx, 3)
	t.UpdateCount(tx)
	return nil
}

// MarkAsFound will mark a specific clue as found for a team
func (t *Team) MarkAsFound(db *gorm.DB, code string) error {
	db.Model(&ClueLog{}).
		Where("clue_code = ? and team = ?", code, t.Code).
		Updates(map[string]interface{}{"found": time.Now(), "active": false}).
		Omit("Clue")
	t.UpdateCount(db)
	t.ActivateClues(db, 3)
	return nil
}

// MarkAsUnfound will remove the found time from a clue for a team
func (t *Team) MarkAsUnfound(db *gorm.DB, code string) error {
	// Update the clue with zero time
	db.Model(&ClueLog{}).
		Where("clue_code = ? and team = ?", code, t.Code).
		Updates(map[string]interface{}{"found": time.Time{}}).
		Omit("Clue")
	t.UpdateCount(db)
	t.ActivateClues(db, 3)
	return nil
}
