package public

import (
	"net/http"
	"strings"

	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/models"
)

// Start begins the game for the team. Prints out their first clue
func Start(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Start"

	r.ParseForm()
	teamCode := strings.ToUpper(r.Form.Get("code"))
	team := models.Team{}

	result := env.DB.Where("code == ?", teamCode).Find(&team)
	if result.RowsAffected != 1 {
		session, _ := env.Session.Get(r, "trace")
		session.AddFlash(flash.Message{Message: "That's not a valid code. Try again."})
		session.Save(r, w)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusBadRequest)
		return nil
	}

	team.Name = r.Form.Get("team")
	team.Start(&env.DB)
	env.DB.Save(&team)
	data["team"] = team

	session, _ := env.Session.Get(r, "trace")
	session.Values["code"] = team.Code
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	data["messages"] = flash.Get(w, r)
	return render(w, data, "start/index.html")
}
