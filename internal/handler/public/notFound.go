package public

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
)

// NotFound is the 404 handler.
func NotFound(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["title"] = "Error 404"
	render(w, data, "errors/404.html")
}

// Error404 is a directly accessible 404 page (via /404)
func Error404(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Error 404"
	session, _ := env.Session.Get(r, "trace")
	data["messages"] = flash.Get(session, w, r)
	return render(w, data, "errors/404notFound.html")
}
