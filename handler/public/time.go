package public

import (
	"fmt"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/handler"
)

// Time returns the amount of time left in the current game.
func Time(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "69 minutes remaining")

	return nil
}
