package public

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
)

// Login handles user logins
func Login(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	session, _ := env.Session.Get(r, "trace")
	data := make(map[string]interface{})
	data["title"] = "Login"
	data["messages"] = flash.Get(session, w, r)

	return render(w, data, "session/login.html")
}

// Register handles user registrations
func Register(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	session, _ := env.Session.Get(r, "trace")
	data := make(map[string]interface{})
	data["title"] = "Register"
	data["messages"] = flash.Get(session, w, r)

	return render(w, data, "session/register.html")
}
