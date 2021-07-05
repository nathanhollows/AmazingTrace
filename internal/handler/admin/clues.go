package admin

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// Clues lists all available clues
func Clues(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	clues := []models.Clue{}
	result := env.DB.Find(&clues)

	data := make(map[string]interface{})
	if result.RowsAffected > 0 {
		data["clues"] = clues
	}
	data["title"] = "Clues | Admin"
	data["messages"] = flash.Get(w, r)

	templates := template.Must(template.ParseFiles(
		"web/templates/admin.html",
		"web/templates/flash.html",
		"web/views/admin/clues/index.html"))

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return nil
}

// CreateClue saves the posted clue
func CreateClue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	if r.Method != http.MethodPost {
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}

	r.ParseForm()
	clue := models.Clue{}
	clue.Location = r.PostFormValue("location")
	clue.Clue = r.PostFormValue("clue")

	result := env.DB.Create(&clue)

	clue.GeneratePoster()

	if result.Error != nil {
		flash.Set(w, r, flash.Message{Message: "Could not save clue", Style: "warning"})
	} else {
		flash.Set(w, r, flash.Message{Message: "New clue added", Style: "success"})
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil
}

// DeleteClue deletes the posted clue
func DeleteClue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	if r.Method != http.MethodPost {
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}

	r.ParseForm()
	clue := models.Clue{}
	result := env.DB.Where("code = ?", r.PostFormValue("code")).Find(&clue)
	if result.Error != nil {
		flash.Set(w, r, flash.Message{Message: "Could not delete clue", Style: "warning"})
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}

	result = env.DB.Where("code = ?", r.PostFormValue("code")).Delete(&clue)
	if result.Error != nil {
		flash.Set(w, r, flash.Message{Message: "Could not delete clue", Style: "warning"})
	} else {
		flash.Set(w, r, flash.Message{Message: "Clue successfully deleted", Style: "success"})
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil
}
