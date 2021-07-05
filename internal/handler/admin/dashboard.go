package admin

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
	"gorm.io/gorm/clause"
)

// Dashboard shows an overview of the game
func Dashboard(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Dashboard"

	teams := []models.Team{}
	env.DB.Where("started == 1").Preload(clause.Associations).Preload("ClueLog.Clue").Order("found desc, updated_at desc").Find(&teams)
	data["teams"] = teams

	session, _ := env.Session.Get(r, "trace")
	data["messages"] = flash.Get(session, w, r)
	return render(w, data, "dashboard/index.html")
}

// FastForward completes three clues for a team.
func FastForward(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	r.ParseForm()
	teamCode := r.PostFormValue("code")

	index, _ := env.Manager.GetTeam(teamCode)
	team := &env.Manager.Teams[index]
	team.FastForward()
	http.Redirect(w, r, "/admin", 301)
	return nil
}

// Hinder completes three clues for a team.
func Hinder(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	r.ParseForm()
	teamCode := r.PostFormValue("code")

	index, _ := env.Manager.GetTeam(teamCode)
	team := &env.Manager.Teams[index]
	team.Hinder()
	http.Redirect(w, r, "/admin", 301)
	return nil
}

// CreateTeam completes three clues for a team.
func CreateTeam(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	env.Manager.CreateTeams(3)
	http.Redirect(w, r, "/admin", 301)
	return nil
}
