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

	session, _ := env.Session.Get(r, "trace")
	data["messages"] = flash.Get(session, w, r)
	return render(w, data, "schedule/index.html")
}

// CreateSchedule saves the new game schedules
func CreateSchedule(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	w.Header().Set("Content-Type", "text/html")

	if r.Method != http.MethodPost {
		return handler.StatusError{http.StatusMethodNotAllowed, errors.New("method not allowed")}
	}

	r.ParseForm()
	schedule := models.Schedule{}
	schedule.SetTimes(r.PostFormValue("date"), r.PostFormValue("start"), r.PostFormValue("end"))

	result := env.DB.Create(&schedule)

	if result.Error != nil {
		return handler.StatusError{http.StatusInternalServerError, errors.New("could not save schedule")}
	}

	schedules := []models.Schedule{}
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
		return handler.StatusError{http.StatusNotFound, errors.New("this schedule does not exist")}
	}

	if r.Method == http.MethodDelete {
		result = env.DB.Delete(&schedule)
		if result.Error != nil {
			return handler.StatusError{http.StatusInternalServerError, errors.New("this schedule could not be deleted")}
		} else {
			return nil
		}
	}

	return handler.StatusError{http.StatusMethodNotAllowed, errors.New("not yet implemented")}
}
