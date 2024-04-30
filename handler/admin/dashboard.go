package admin

import (
	"net/http"
	"time"

	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/models"
	"gorm.io/gorm/clause"
)

// Dashboard shows an overview of the game
func Dashboard(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	env.Data["title"] = "Dashboard"

	var codeCount int64
	env.DB.Model(&models.Team{}).Count(&codeCount)
	env.Data["code_count"] = int(codeCount)
	if codeCount != 0 {
		teams := []models.Team{}
		env.DB.Where("started == 1").Preload(clause.Associations).Preload("Clues.Clue").Order("found desc, updated_at asc").Find(&teams)
		env.Data["teams"] = teams
	}

	// Find the game with the next *end* time, in case we're in the middle of a game
	game := models.Game{}
	result := env.DB.Where("end_time > ?", time.Now()).Order("end_time ASC").Limit(1).Find(&game)
	if result.RowsAffected != 0 {
		env.Data["game"] = game
	}

	var clues []models.Clue
	env.DB.Find(&clues)
	env.Data["clues"] = clues
	env.Data["clue_count"] = len(clues)

	env.Data["messages"] = flash.Get(w, r)
	return render(w, env.Data, "dashboard/index.html")
}

// DashboardTable renders only the table part of the dashboard
func DashboardTable(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	teams := []models.Team{}
	env.DB.Where("started == 1").Preload(clause.Associations).Preload("Clues.Clue").Order("found desc, updated_at asc").Find(&teams)
	env.Data["teams"] = teams

	return renderFragment(w, env.Data, "dashboard/table.html")
}
