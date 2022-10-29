package public

import (
	"net/http"
	"time"

	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/models"
)

// Clues shows a team all of the clues they have unlocked
func Clues(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	session, err := env.Session.Get(r, "trace")
	if err != nil {
		session.AddFlash(flash.Message{Message: "Something went wrong. You don't seem to be registered.", Style: "warning"})
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	teamCode := session.Values["code"]
	if teamCode == nil {
		session.AddFlash(flash.Message{Message: "You need to be logged in to play.", Style: "warning"})
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil
	}

	data := make(map[string]interface{})
	data["title"] = "Clues"

	team := &models.Team{}
	env.DB.Where("code = ?", teamCode).Preload("Clues.Clue").Find(&team)
	data["team"] = team

	var solved int64
	env.DB.Model(&models.ClueLog{}).
		Where("team = ? AND found <> ?", team.Code, time.Time{}).Count(&solved)
	data["solved"] = solved

	data["messages"] = flash.Get(w, r)
	session.Save(r, w)
	return render(w, data, "clues/index.html")
}
