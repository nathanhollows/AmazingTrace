package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/models"
	"gorm.io/gorm/clause"
)

// Teams shows all teams and their status
func Teams(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Team Codes"
	data["messages"] = flash.Get(w, r)

	teams := []models.Team{}
	env.DB.Where("not started").Model(&models.Team{}).Find(&teams)
	data["teams"] = teams

	return render(w, data, "teams/index.html")
}

// GenerateTeams will generate X numbers of teams
func GenerateTeams(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	session, _ := env.Session.Get(r, "trace")

	r.ParseForm()
	i := r.PostFormValue("count")
	count, err := strconv.Atoi(i)
	if err != nil {
		session.AddFlash(flash.Message{Message: "Could not generate teams", Style: "warning"})
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		return nil
	}

	if count == 0 {
		env.DB.Unscoped().Where("started <> 1").Delete(&models.Team{})
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		return nil
	}

	for i := 0; i < count; i++ {
		result := env.DB.Model(&models.Team{}).Create(&models.Team{})
		if result.Error != nil {
			if !strings.Contains(result.Error.Error(), "UNIQUE constraint failed:") {
				session.AddFlash(flash.Message{Message: "Something went wrong generating the teams", Style: "warning"})
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
				return nil
			}
			i-- // Try again
		}
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	return nil

}

// TeamInspect shows the current status of a team
func TeamInspect(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Team Codes"

	code := chi.URLParam(r, "code")

	team := models.Team{}
	env.DB.Model(&team).Where("code = ?", code).Preload("Clues.Clue").Preload(clause.Associations).Find(&team)
	data["team"] = team

	return renderFragment(w, data, "teams/inspect.html")
}

// FastForward completes three clues for a team.
func FastForward(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	code := chi.URLParam(r, "code")

	team := &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	if team.Code == "" {
		http.Error(w, "team not found", http.StatusNotFound)
		return nil
	}

	team.FastForward(&env.DB, 1)
	team = &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	data["team"] = team

	return renderFragment(w, data, "teams/scans.html")
}

// Shuffle randomises the three clues a team is looking for.
func Shuffle(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	code := chi.URLParam(r, "code")

	team := &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	if team.Code == "" {
		http.Error(w, "team not found", http.StatusNotFound)
		return nil
	}
	team.Shuffle(&env.DB)
	team = &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	data["team"] = team

	return renderFragment(w, data, "teams/scans.html")
}

// Rewind swaps the last sovled clue back into the game.
func Rewind(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	code := chi.URLParam(r, "code")

	team := &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	if team.Code == "" {
		http.Error(w, "team not found", http.StatusNotFound)
		return nil
	}
	team.Rewind(&env.DB)
	team = &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	data["team"] = team

	return renderFragment(w, data, "teams/scans.html")
}

// Given a clue code, mark it as solved for a team.
func MarkAsFound(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	code := chi.URLParam(r, "team")
	clueCode := chi.URLParam(r, "clue")

	team := &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	if team.Code == "" {
		http.Error(w, "team not found", http.StatusNotFound)
		return nil
	}
	team.MarkAsFound(&env.DB, clueCode)
	team = &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	data["team"] = team

	return renderFragment(w, data, "teams/scans.html")
}

// Given a clue code, mark it as unsolved for a team.
func MarkAsUnfound(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	code := chi.URLParam(r, "team")
	clueCode := chi.URLParam(r, "clue")

	team := &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	if team.Code == "" {
		http.Error(w, "team not found", http.StatusNotFound)
		return nil
	}
	team.MarkAsUnfound(&env.DB, clueCode)
	team = &models.Team{}
	env.DB.Where("code = ?", code).Preload("Clues.Clue").Find(&team)
	data["team"] = team

	return renderFragment(w, data, "teams/scans.html")
}
