package admin

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// Teams shows all teams and their status
func Teams(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	data := make(map[string]interface{})
	data["title"] = "Team Codes | Admin"
	data["messages"] = flash.Get(w, r)

	teams := []models.Team{}
	env.DB.Where("not started").Model(&models.Team{}).Find(&teams)
	data["teams"] = teams

	templates := template.Must(template.ParseFiles(
		"web/templates/admin.html",
		"web/templates/flash.html",
		"web/views/admin/teams/index.html"))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}

// GenerateTeams will generate X numbers of teams
func GenerateTeams(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	r.ParseForm()
	i := r.PostFormValue("count")
	count, err := strconv.Atoi(i)
	if err != nil {
		flash.Set(w, r, flash.Message{Message: "Could not generate teams", Style: "warning"})
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
				flash.Set(w, r, flash.Message{Message: "Something went wrong generating the teams", Style: "warning"})
				http.Redirect(w, r, r.Header.Get("Referer"), 302)
				return nil
			}
			i-- // Try again
		}
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil

}
