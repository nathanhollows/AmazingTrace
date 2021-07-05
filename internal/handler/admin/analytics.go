package admin

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
)

// Analytics shows player stats
func Analytics(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	session, _ := env.Session.Get(r, "trace")
	data["title"] = "Analytics"
	data["messages"] = flash.Get(session, w, r)
	return render(w, data, "analytics/index.html")
}
