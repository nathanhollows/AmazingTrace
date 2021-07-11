package admin

import (
	"net/http"
	"time"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
	"gorm.io/gorm/clause"
)

// Dashboard shows an overview of the game
func Dashboard(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})
	data["title"] = "Dashboard"

	var code_count int64
	env.DB.Model(&models.Team{}).Count(&code_count)
	data["code_count"] = int(code_count)
	if code_count != 0 {
		teams := []models.Team{}
		env.DB.Where("started == 1").Preload(clause.Associations).Preload("ClueLog.Clue").Order("found desc, updated_at desc").Find(&teams)
		data["teams"] = teams
	}

	// Find the game with the next *end* time, in case we're in the middle of a game
	game := models.Game{}
	result := env.DB.Where("end_time > ?", time.Now()).Order("end_time ASC").First(&game)
	if result.Error == nil {
		data["game"] = game
	}

	var clue_count int64
	env.DB.Model(&models.Clue{}).Count(&clue_count)
	data["clue_count"] = clue_count

	session, _ := env.Session.Get(r, "trace")
	data["messages"] = flash.Get(session, w, r)
	return render(w, data, "dashboard/index.html")
}
