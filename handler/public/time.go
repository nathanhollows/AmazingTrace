package public

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/models"
	"github.com/xeonx/timeago"
)

// Time returns the amount of time left in the current game.
func Time(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	schedules := []models.Schedule{}
	env.DB.Order("start").Find(&schedules)

	timeConf := timeago.NoMax(timeago.English)

	// Find the current
	for _, s := range schedules {
		if s.Start.Before(time.Now()) && s.End.After(time.Now()) {
			s := timeago.Config(timeConf).Format(s.End)
			fmt.Fprintln(w, "Game ends "+s)
			return nil
		}
	}

	// Check if there is a schedule in the future
	for _, s := range schedules {
		if s.Start.After(time.Now()) {
			fmt.Fprint(w, "Game starts "+timeago.Config(timeConf).Format(s.Start))
			return nil
		}
	}

	fmt.Fprint(w, "No upcoming games <a href=\"/admin/schedule\">scheduled</a>.")
	return nil
}
