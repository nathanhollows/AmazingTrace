package public

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// Index is the homepage of the game.
// Prints a very simple page asking only for a team code.
func Index(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	session, _ := env.Session.Get(r, "trace")
	data := make(map[string]interface{})
	data["messages"] = flash.Get(session, w, r)

	// Check if team exists
	teamCode := session.Values["code"]
	result := env.DB.Where("Code = ?", teamCode).Find(&models.Team{})
	if result.RowsAffected == 1 {
		http.Redirect(w, r, "/clues", 302)
	}

	return render(w, data, "index/index.html")
}
