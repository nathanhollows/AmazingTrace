package admin

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
	"gorm.io/gorm/clause"
)

// Dashboard shows an overview of the game
func Dashboard(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	data := make(map[string]interface{})
	data["title"] = "Dashboard | Admin"

	teams := []models.Team{}
	env.DB.Where("started == 1").Preload(clause.Associations).Preload("ClueLog.Clue").Order("found desc, updated_at desc").Find(&teams)
	data["teams"] = teams

	data["messages"] = flash.Get(w, r)

	fm := template.FuncMap{
		"divide": func(a, b int) float32 {
			if a == 0 || b == 0 {
				return 0
			}
			return float32(a) / float32(b)
		},
		"progress": func(a, b int) float32 {
			if a == 0 || b == 0 {
				return 0
			}
			return float32(a) / float32(b) * 100
		},
		"add": func(a, b int) int {
			return a + b
		},
	}

	templates := template.Must(template.New("dashboard").Funcs(fm).ParseFiles(
		"web/templates/admin.html",
		"web/templates/flash.html",
		"web/views/admin/dashboard/index.html"))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
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
