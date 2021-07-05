package admin

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/handler"
)

// Analytics shows player stats
func Analytics(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	data := make(map[string]interface{})
	data["title"] = "Analytics | Admin"

	templates := template.Must(template.ParseFiles(
		"web/templates/admin.html",
		"web/templates/flash.html",
		"web/views/admin/analytics/index.html"))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}
