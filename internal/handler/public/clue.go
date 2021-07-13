package public

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// Clue handles the scanned URL.
func Clue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	session, _ := env.Session.Get(r, "trace")

	// Check if any team is currently logged in
	team, err := models.Team{}.Get(&env.DB, session.Values["code"])
	if err != nil {
		session.AddFlash(flash.Message{Message: "You need to be logged in to scan clues.", Style: "warning"})
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
		return nil
	}

	// Check if the clue is one the team needed to find
	err = team.CheckClue(&env.DB, chi.URLParam(r, "code"))
	if err != nil {
		session.AddFlash(flash.Message{Message: err.Error(), Style: "warning"})
		session.Save(r, w)
		http.Redirect(w, r, "/clues", 302)
		return nil
	}

	session.AddFlash(flash.Message{Message: "👏👏👏 Congratulations! 👏👏👏", Style: "success text-center"})
	session.Save(r, w)
	http.Redirect(w, r, "/clues", 302)
	return nil
}
