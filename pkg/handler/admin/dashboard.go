package admin

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/pkg/handler"
)

// Dashboard shows an overview of the game
func Dashboard(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	data := make(map[string]interface{})
	data["title"] = "Admin"

	templates := template.Must(template.ParseFiles(
		"../web/templates/admin.html",
		"../web/templates/flash.html",
		"../web/views/admin/dashboard/index.html"))

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
