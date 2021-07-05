package public

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
)

// Login handles user logins
func Login(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	templates := template.Must(template.ParseFiles(
		"web/templates/index.html",
		"web/views/session/login.html"))

	if err := templates.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}

// Register handles user registrations
func Register(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	switch r.Method {
	case http.MethodPost:
	}

	env.Data["messages"] = flash.Get(w, r)

	templates := template.Must(template.ParseFiles(
		"web/templates/index.html",
		"web/views/session/register.html"))

	if err := templates.ExecuteTemplate(w, "base", env.Data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}
