package admin

import (
	"net/http"

	"github.com/nathanhollows/AmazingTrace/internal/flash"
	"github.com/nathanhollows/AmazingTrace/internal/handler"
	"github.com/nathanhollows/AmazingTrace/internal/models"
)

// Clues lists all available clues
func Clues(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	clues := []models.Clue{}
	result := env.DB.Find(&clues)

	data := make(map[string]interface{})
	if result.RowsAffected > 0 {
		data["clues"] = clues
	}
	data["title"] = "Clues"
	session, _ := env.Session.Get(r, "trace")
	data["messages"] = flash.Get(session, w, r)
	return render(w, data, "clues/index.html")
}

// CreateClue saves the posted clue
func CreateClue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	session, _ := env.Session.Get(r, "trace")

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
		session.AddFlash(flash.Message{Message: "Could not save clue", Style: "warning"})
		session.Save(r, w)
	} else {
		session.AddFlash(flash.Message{Message: "New clue added", Style: "success"})
		session.Save(r, w)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil
}

// DeleteClue deletes the posted clue
func DeleteClue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	session, _ := env.Session.Get(r, "trace")

	if r.Method != http.MethodPost {
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}

	r.ParseForm()
	clue := models.Clue{}
	result := env.DB.Where("code = ?", r.PostFormValue("code")).Find(&clue)
	if result.Error != nil {
		session.AddFlash(flash.Message{Message: "Could not delete clue", Style: "warning"})
		session.Save(r, w)
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}

	result = env.DB.Where("code = ?", r.PostFormValue("code")).Delete(&clue)
	if result.Error != nil {
		session.AddFlash(flash.Message{Message: "Could not delete clue", Style: "warning"})
		session.Save(r, w)
	} else {
		session.AddFlash(flash.Message{Message: "Clue successfully deleted", Style: "success"})
		session.Save(r, w)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil
}
