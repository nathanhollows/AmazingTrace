package admin

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
)

// Schedule lets the game admin set game variables
func Schedule(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Schedule"
	session, _ := env.Session.Get(r, "trace")
	data["messages"] = flash.Get(session, w, r)
	return render(w, data, "schedule/index.html")
}
