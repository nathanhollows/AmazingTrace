package public

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/helpers"
	"github.com/nathanhollows/AmazingTrace/models"
)

// Login handles user logins
func Login(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	env.Data["title"] = "Login"

	r.ParseForm()

	if r.Method == http.MethodPost {
		username := r.Form.Get("username")

		var user models.Admin
		env.DB.Model(&user).Where("username = ?", username).Find(&user)

		if user.CheckHashPassword(r.Form.Get("password")) {
			session, err := env.Session.Get(r, "admin")
			if err != nil || session.Values["id"] == nil {
				session, _ = env.Session.New(r, "admin")
				session.Options.HttpOnly = true
				session.Options.SameSite = http.SameSiteStrictMode
				session.Options.Secure = true
				id := uuid.New()
				session.Values["id"] = id.String()
				session.Save(r, w)
				http.Redirect(w, r, helpers.URL("admin"), http.StatusFound)
				return nil
			}
		} else {
			flash.Error(w, r, "Invalid username or password")
		}
	}

	env.Data["messages"] = flash.Get(w, r)
	return render(w, env.Data, "session/login.html")
}

// Register handles user registrations
func Register(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	env.Data["title"] = "Register"
	env.Data["messages"] = flash.Get(w, r)

	return render(w, env.Data, "session/register.html")
}
