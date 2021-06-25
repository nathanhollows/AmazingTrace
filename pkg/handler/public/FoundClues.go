package public

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/pkg/flash"
	"github.com/nathanhollows/AmazingTrace/pkg/handler"
	"github.com/nathanhollows/AmazingTrace/pkg/models"
)

// FoundClues shows a team all of the clues they have unlocked
func FoundClues(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	data := make(map[string]interface{})
	data["title"] = "Clues | The Amazing Trace"
	data["messages"] = flash.Get(w, r)

	r.ParseForm()
	teamCode := r.PostForm.Get("code")
	team := &models.Team{}
	env.DB.Where("code = ?", teamCode).Find(&team)

	data["team"] = team

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
