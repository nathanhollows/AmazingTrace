package public

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/nathanhollows/AmazingTrace/pkg/game"
	"github.com/nathanhollows/AmazingTrace/pkg/handler"
)

// Clue handles the scanned URL.
func Clue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	type Data struct {
		Clue     game.Clue
		Team     game.Team
		TimeLeft float64
	}
	data := Data{}

	var st = "👏👏👏 Congratulations! 👏👏👏"
	fmt.Print(st)

	// TODO: Make this variable.
	end := time.Date(2021, time.February, 24, 12, 0, 0, 0, time.Local)
	data.TimeLeft = math.Floor(end.Sub(time.Now()).Minutes())

	if len(r.URL.String()) != 6 {
		templates := template.Must(template.ParseFiles(
			"../web/templates/index.html",
			"../web/templates/flash.html",
			"../web/views/errors/notfound.html"))

		if err := templates.ExecuteTemplate(w, "base", env.Data); err != nil {
			http.Error(w, err.Error(), 0)
			log.Print("Template executing error: ", err)
		}

		return nil
	}

	clueCode := (r.URL.String()[1:6])
	clueCode = strings.ToUpper(clueCode)
	var page string = "clue/404"

	clue, err := env.Manager.GetClue(clueCode)
	if err != nil {
		page = "clue/404"
	}
	data.Clue = clue

	// TODO: Create a seperate handler for POST.
	if r.Method == "POST" {
		r.ParseForm()
		teamCode := r.Form.Get("code")
		index, err := env.Manager.GetTeam(teamCode)
		team := &game.Team{}
		if err != nil {
			page = "clue/notateam"
		} else {
			team = &env.Manager.Teams[index]
			if team.Solve(clueCode) != nil {
				page = "clues/index"
				team.Status = "error"
			} else {
				page = "clues/index"
				team.Status = "success"
			}
		}
		team.CheckIn()
		data.Team = *team
	} else {
		session, err := env.Session.Get(r, "trace")
		code := session.Values["code"]
		teamCode := fmt.Sprintf("%v", code)
		index, err := env.Manager.GetTeam(teamCode)
		if err != nil {
			page = "clue/index"
		} else {
			team := &env.Manager.Teams[index]
			if err != nil {
				page = "clue/notateam"
			} else {
				team = &env.Manager.Teams[index]
				page = "clue/index"
			}
			team.CheckIn()
			data.Team = *team
		}
	}

	_, err = env.Manager.GetClue(clueCode)
	if err != nil {
		page = "clue/404"
	}

	templates := template.Must(template.ParseFiles(
		"../web/templates/index.html",
		"../web/views/"+page+".html"))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}
