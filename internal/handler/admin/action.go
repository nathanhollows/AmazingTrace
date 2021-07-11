package admin

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/handler"
)

// FastForward completes three clues for a team.
func FastForward(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	r.ParseForm()
	teamCode := r.PostFormValue("code")

	index, _ := env.Manager.GetTeam(teamCode)
	team := &env.Manager.Teams[index]
	team.FastForward()
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
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
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
	return nil
}
