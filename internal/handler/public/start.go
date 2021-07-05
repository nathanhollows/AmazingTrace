package public

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// Start begins the game for the team. Prints out their first clue
func Start(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	data := make(map[string]interface{})
	data["title"] = "Start | The Amazing Trace"

	r.ParseForm()
	teamCode := strings.ToUpper(r.Form.Get("code"))
	team := models.Team{}

	result := env.DB.Where("code == ?", teamCode).Find(&team)
	if result.RowsAffected != 1 {
		session, _ := env.Session.Get(r, "trace")
		session.AddFlash(flash.Message{Message: "That's not a valid code. Try again.", Style: "warning"})
		session.Save(r, w)
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return nil
	}

	team.Start(&env.DB)
	env.DB.Save(&team)
	data["team"] = team

	session, _ := env.Session.Get(r, "trace")
	session.Values["code"] = team.Code
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	templates := template.Must(template.ParseFiles(
		"web/templates/index.html",
		"web/templates/flash.html",
		"web/views/public/start/index.html"))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}
