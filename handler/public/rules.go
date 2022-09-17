package public

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
)

// Rules shows the player the 'start' page rules, without any of the setup logic executing
func Rules(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Rules"
	session, _ := env.Session.Get(r, "trace")
	data["messages"] = flash.Get(session, w, r)
	return render(w, data, "rules/index.html")
}
