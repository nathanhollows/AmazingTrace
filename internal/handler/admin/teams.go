package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// Teams shows all teams and their status
func Teams(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Team Codes"
	session, _ := env.Session.Get(r, "trace")
	data["messages"] = flash.Get(session, w, r)

	teams := []models.Team{}
	env.DB.Where("not started").Model(&models.Team{}).Find(&teams)
	data["teams"] = teams

	return render(w, data, "teams/index.html")
}

// GenerateTeams will generate X numbers of teams
func GenerateTeams(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	session, _ := env.Session.Get(r, "trace")

	r.ParseForm()
	i := r.PostFormValue("count")
	count, err := strconv.Atoi(i)
	if err != nil {
		session.AddFlash(flash.Message{Message: "Could not generate teams", Style: "warning"})
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return nil
	}

	if count == 0 {
		env.DB.Unscoped().Where("started <> 1").Delete(&models.Team{})
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
		return nil
	}

	for i := 0; i < count; i++ {
		result := env.DB.Model(&models.Team{}).Create(&models.Team{})
		if result.Error != nil {
			if !strings.Contains(result.Error.Error(), "UNIQUE constraint failed:") {
				session.AddFlash(flash.Message{Message: "Something went wrong generating the teams", Style: "warning"})
				http.Redirect(w, r, r.Header.Get("Referer"), 302)
				return nil
			}
			i-- // Try again
		}
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil

}
