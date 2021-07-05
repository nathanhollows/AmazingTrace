package public

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// Index is the homepage of the game.
// Prints a very simple page asking only for a team code.
func Index(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	session, _ := env.Session.Get(r, "trace")
	data := make(map[string]interface{})

	// Check if team exists
	teamCode := session.Values["code"]
	result := env.DB.Where("Code = ?", teamCode).Find(&models.Team{})
	if result.RowsAffected == 1 {
		http.Redirect(w, r, "/clues", 302)
	}
	data["title"] = "The Amazing Trace"

	// Get flashed messages
	data["messages"] = session.Flashes()
	session.Save(r, w)

	templates := template.Must(template.ParseFiles(
		"web/templates/index.html",
		"web/templates/flash.html",
		"web/views/public/index/index.html"))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}
