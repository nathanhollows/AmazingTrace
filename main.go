// AmazingTrace is a QR code based scavenger hunt
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/nathanhollows/AmazingTrace/internal/filesystem"
	"github.com/nathanhollows/AmazingTrace/internal/game"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/handler/admin"
	"github.com/nathanhollows/AmazingTrace/internal/handler/public"
	"github.com/nathanhollows/AmazingTrace/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var router *chi.Mux
var env handler.Env

func init() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(5))

	var store = sessions.NewCookieStore([]byte("trace"))
	store.Options.SameSite = http.SameSiteStrictMode

	db, err := gorm.Open(sqlite.Open("trace.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	env = handler.Env{
		Manager: game.Manager{},
		Session: store,
		DB:      *db,
		Data:    make(map[string]interface{}),
	}
}

func main() {
	env.DB.AutoMigrate(
		&models.Area{},
		&models.Clue{},
		&models.User{},
		&models.ClueLog{},
		&models.Team{},
	)
	routes()
	fmt.Println(http.ListenAndServe(":8000", router))
}

// Set up the routes needed for the game.
func routes() {
	router.Handle("/", handler.Handler{Env: &env, H: public.Index})

	router.Handle("/start", handler.Handler{Env: &env, H: public.Start})
	router.Handle("/clues", handler.Handler{Env: &env, H: public.FoundClues})
	router.Handle("/rules", handler.Handler{Env: &env, H: public.Rules})
	router.Handle("/{code:[A-z]{4}}", handler.Handler{Env: &env, H: public.Clue})

	router.Handle("/login", handler.Handler{Env: &env, H: public.Login})
	router.Handle("/register", handler.Handler{Env: &env, H: public.Register})

	router.Handle("/admin", handler.Admin{Env: &env, H: admin.Dashboard})
	router.Handle("/admin/ff", handler.Admin{Env: &env, H: admin.FastForward})
	router.Handle("/admin/hinder", handler.Admin{Env: &env, H: admin.Hinder})
	router.Handle("/admin/teams", handler.Admin{Env: &env, H: admin.Teams})
	router.Handle("/admin/teams/generate", handler.Admin{Env: &env, H: admin.GenerateTeams})
	router.Handle("/admin/clues", handler.Admin{Env: &env, H: admin.Clues})
	router.Handle("/admin/clues/create", handler.Admin{Env: &env, H: admin.CreateClue})
	router.Handle("/admin/clues/delete", handler.Admin{Env: &env, H: admin.DeleteClue})
	router.Handle("/admin/analytics", handler.Admin{Env: &env, H: admin.Analytics})

	router.Handle("/404", handler.Handler{Env: &env, H: public.Error404})
	router.NotFound(public.NotFound)

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "web/static"))}
	filesystem.FileServer(router, "/public", filesDir)
}
