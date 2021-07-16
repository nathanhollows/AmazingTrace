package admin

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// FastForward completes three clues for a team.
func FastForward(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	r.ParseForm()

	teamCode := r.PostFormValue("code")
	team := &models.Team{}
	env.DB.Where("code = ?", teamCode).Find(&team)
	if team.Code == "" {
		flash.Set(w, r, flash.Message{Style: "danger", Message: "Could not fast forward team. The team could not be found."})
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return nil
	}
	team.FastForward(&env.DB, 1)

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil
}

// Shuffle randomises the three clues a team is looking for.
func Shuffle(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	teamCode := r.PostFormValue("code")
	team := &models.Team{}
	env.DB.Where("code = ?", teamCode).Find(&team)
	if team.Code == "" {
		flash.Set(w, r, flash.Message{Style: "danger", Message: "Could not fast forward team. The team could not be found."})
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return nil
	}
	team.Shuffle(&env.DB)

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil
}

// Rewind swaps the last sovled clue back into the game.
func Rewind(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	teamCode := r.PostFormValue("code")
	team := &models.Team{}
	env.DB.Where("code = ?", teamCode).Find(&team)
	if team.Code == "" {
		flash.Set(w, r, flash.Message{Style: "danger", Message: "Could not fast forward team. The team could not be found."})
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return nil
	}
	team.Rewind(&env.DB)

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil
}
