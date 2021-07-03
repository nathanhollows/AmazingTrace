package public

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/nathanhollows/AmazingTrace/pkg/flash"
	"github.com/nathanhollows/AmazingTrace/pkg/handler"
	"github.com/nathanhollows/AmazingTrace/pkg/models"
)

// FoundClues shows a team all of the clues they have unlocked
func FoundClues(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	session, err := env.Session.Get(r, "trace")
	if err != nil {
		flash.Set(w, r, flash.Message{Message: "Something went wrong. You don't seem to be registered.", Style: "warning"})
		http.Redirect(w, r, "/", 302)
	}
	teamCode := session.Values["code"]
	if teamCode == nil {
		session.AddFlash(flash.Message{Message: "You need to be logged in to play.", Style: "warning"})
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
		return nil
	}

	data := make(map[string]interface{})
	data["title"] = "Clues | The Amazing Trace"

	team := &models.Team{}
	env.DB.Where("code = ?", teamCode).Preload("ClueLog.Clue").Find(&team)
	data["team"] = team

	var solved int64
	env.DB.Model(&models.ClueLog{}).
		Where("team = ? AND found <> ?", team.Code, time.Time{}).Count(&solved)
	data["solved"] = solved
	data["messages"] = session.Flashes()
	session.Save(r, w)

	templates := template.Must(template.ParseFiles(
		"../web/templates/index.html",
		"../web/templates/flash.html",
		"../web/views/public/clues/index.html",
	))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}
