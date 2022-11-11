// AmazingTrace is a QR code based scavenger hunt
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/nathanhollows/AmazingTrace/filesystem"
	"github.com/nathanhollows/AmazingTrace/game"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/handler/admin"
	"github.com/nathanhollows/AmazingTrace/handler/public"
	"github.com/nathanhollows/AmazingTrace/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var router *chi.Mux
var env handler.Env

func init() {
	var store = sessions.NewCookieStore([]byte("trace"))
	store.Options.SameSite = http.SameSiteStrictMode

	db, err := gorm.Open(sqlite.Open("trace.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	data := make(map[string]interface{})

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	data["MAPBOX_KEY"] = os.Getenv("MAPBOX_KEY")

	env = handler.Env{
		Manager: game.Manager{},
		Session: store,
		DB:      *db,
		Data:    data,
	}

	// Set up Admin user
	_, err = models.CreateAdmin(env.DB, os.Getenv("ADMIN_USERNAME"), os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		panic(err)
	}

	// Clean up posters
	go models.CleanUpPosters(db)
}

func main() {
	env.DB.AutoMigrate(
		&models.Admin{},
		&models.Area{},
		&models.Clue{},
		&models.ClueLog{},
		&models.Game{},
		&models.Schedule{},
		&models.Team{},
		&models.User{},
	)
	routes()
	fmt.Println(http.ListenAndServe(":8000", router))
}

// Set up the routes needed for the game.
func routes() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(5))

	router.Handle("/", handler.HandlePublic{Env: &env, H: public.Index})

	router.Handle("/start", handler.HandlePublic{Env: &env, H: public.Start})
	router.Handle("/clues", handler.HandlePublic{Env: &env, H: public.Clues})
	router.Handle("/rules", handler.HandlePublic{Env: &env, H: public.Rules})
	router.Handle("/time", handler.HandlePublic{Env: &env, H: public.Time})
	router.Handle("/{code:[A-z]{4}}", handler.HandlePublic{Env: &env, H: public.Scan})

	router.Handle("/login", handler.HandlePublic{Env: &env, H: public.Login})
	// router.Handle("/register", handler.HandlePublic{Env: &env, H: public.Register})
	// Dashboard
	router.Handle("/admin", handler.HandleAdmin{Env: &env, H: admin.Dashboard})
	router.Handle("/admin/dashboard", handler.HandleAdmin{Env: &env, H: admin.Dashboard})
	router.Handle("/admin/dashboard/table", handler.HandleAdmin{Env: &env, H: admin.DashboardTable})
	// Teams
	router.Handle("/admin/teams", handler.HandleAdmin{Env: &env, H: admin.Teams})
	router.Handle("/admin/teams/inspect/{code:[A-z]{4}}", handler.HandleAdmin{Env: &env, H: admin.TeamInspect})
	router.Handle("/admin/teams/generate", handler.HandleAdmin{Env: &env, H: admin.GenerateTeams})
	router.Handle("/admin/teams/fastforward/{code:[A-z]{4}}", handler.HandleAdmin{Env: &env, H: admin.FastForward})
	router.Handle("/admin/teams/shuffle/{code:[A-z]{4}}", handler.HandleAdmin{Env: &env, H: admin.Shuffle})
	router.Handle("/admin/teams/rewind/{code:[A-z]{4}}", handler.HandleAdmin{Env: &env, H: admin.Rewind})
	router.Handle("/admin/teams/solve/{team:[A-z]{4}}/{clue:[A-z]{4}}", handler.HandleAdmin{Env: &env, H: admin.MarkAsFound})
	router.Handle("/admin/teams/unsolve/{team:[A-z]{4}}/{clue:[A-z]{4}}", handler.HandleAdmin{Env: &env, H: admin.MarkAsUnfound})

	// Clues
	router.Handle("/admin/clues", handler.HandleAdmin{Env: &env, H: admin.Clues})
	router.Handle("/admin/clues/{code:[A-z]{4}}", handler.HandleAdmin{Env: &env, H: admin.ChangeClues})
	router.Handle("/admin/clues/create", handler.HandleAdmin{Env: &env, H: admin.CreateClue})
	router.Handle("/admin/clues/new", handler.HandleAdmin{Env: &env, H: admin.NewClue})
	router.Handle("/admin/clues/edit/{code:[A-z]{4}}", handler.HandleAdmin{Env: &env, H: admin.EditClue})

	// Schedule
	router.Handle("/admin/schedule", handler.HandleAdmin{Env: &env, H: admin.Schedule})
	router.Handle("/admin/schedule/create", handler.HandleAdmin{Env: &env, H: admin.CreateSchedule})
	router.Handle("/admin/schedule/table", handler.HandleAdmin{Env: &env, H: admin.ScheduleTable})
	router.Handle("/admin/schedule/{ID}", handler.HandleAdmin{Env: &env, H: admin.ChangeSchedule})

	router.Handle("/404", handler.HandlePublic{Env: &env, H: public.Error404})
	router.NotFound(public.NotFound)

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "assets"))}
	filesystem.FileServer(router, "/static", filesDir)
}
