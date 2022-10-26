package admin

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/models"
)

// Schedule lets the game admin set game variables
func Schedule(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Schedule"

	schedule := []models.Schedule{}
	env.DB.Model(&schedule).Find(&schedule)
	data["schedule"] = schedule

	data["messages"] = flash.Get(w, r)
	return render(w, data, "schedule/index.html")
}

// CreateSchedule saves the new game schedules
func CreateSchedule(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	w.Header().Set("Content-Type", "text/html")
	schedules := []models.Schedule{}
	env.DB.Find(&schedules)
	data["schedule"] = schedules

	if r.Method != http.MethodPost {
		return handler.StatusError{Code: http.StatusMethodNotAllowed, Err: errors.New("method not allowed")}
	}

	r.ParseForm()
	schedule := models.Schedule{}
	schedule.SetTimes(r.PostFormValue("date"), r.PostFormValue("start"), r.PostFormValue("end"))

	// Make sure the start time is before the end time
	if schedule.Start.After(schedule.End) {
		data["error"] = "The start time must be before the end time"
		return renderFragment(w, data, "schedule/table.html")
	} else if schedule.Start.Equal(schedule.End) {
		data["error"] = "The start time must be before the end time"
		return renderFragment(w, data, "schedule/table.html")
	}

	// Make sure the schedule doesn't overlap with any other schedules
	existing := []models.Schedule{}
	result := env.DB.Find(&existing)
	if result.Error != nil {
		data["error"] = result.Error.Error()
		return renderFragment(w, data, "schedule/table.html")
	}
	// Check if there are any overlapping schedules
	for _, s := range existing {
		if schedule.Overlaps(s) {
			data["error"] = "The schedule overlaps with another schedule"
			return renderFragment(w, data, "schedule/table.html")
		}
	}

	result = env.DB.Create(&schedule)

	if result.Error != nil {
		return handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("could not save schedule")}
	}

	env.DB.Find(&schedules)
	data["schedule"] = schedules
	return renderFragment(w, data, "schedule/table.html")
}

// ScheduleTable renders just the table fragment for this page
func ScheduleTable(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})

	schedule := []models.Schedule{}
	env.DB.Find(&schedule)
	data["schedule"] = schedule

	return renderFragment(w, data, "schedule/table.html")
}

// ChangeSchedule updates a given schedule
func ChangeSchedule(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "ID")

	schedule := models.Schedule{}
	result := env.DB.Where("id = ?", id).Find(&schedule)
	if result.Error != nil {
		return handler.StatusError{Code: http.StatusNotFound, Err: errors.New("this schedule does not exist")}
	}

	if r.Method == http.MethodDelete {
		result = env.DB.Delete(&schedule)
		if result.Error != nil {
			return handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("this schedule could not be deleted")}
		} else {
			return nil
		}
	}

	return handler.StatusError{Code: http.StatusMethodNotAllowed, Err: errors.New("not yet implemented")}
}
