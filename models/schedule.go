package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Schedule keeps track of different game times
type Schedule struct {
	gorm.Model
	Date  time.Time
	Start time.Time
	End   time.Time
}

func (s *Schedule) SetTimes(date string, start string, end string) error {
	var err error
	s.Date, err = time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}
	s.Start, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprint(date, " ", start), time.Local)
	if err != nil {
		return err
	}
	s.End, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprint(date, " ", end), time.Local)
	if err != nil {
		return err
	}
	return nil
}

func (s *Schedule) Overlaps(other Schedule) bool {
	return s.Start.Before(other.End) && s.End.After(other.Start)
}
