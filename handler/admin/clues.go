package admin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/models"
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
	w.Header().Set("Content-Type", "text/html")

	if r.Method != http.MethodPost {
		return handler.StatusError{http.StatusMethodNotAllowed, errors.New("method not allowed")}
	}

	r.ParseForm()
	clue := models.Clue{}
	clue.Location = r.PostFormValue("location")
	clue.Clue = r.PostFormValue("clue")

	result := env.DB.Create(&clue)

	clue.GeneratePoster()

	if result.Error != nil {
		return handler.StatusError{http.StatusInternalServerError, errors.New("could not save clue")}
	}

	fmt.Fprintf(w, `
<tr>
	<th scope="row">%v</th>
	<td>%v</td>
	<td>%v</td>
	<td>                   
		<a href="/static/img/posters/%v.png">Get Poster</a>
	</td>
	<td>
		<a name="code" hx-delete="/admin/clues/%v">Delete</a> 
	</td>
</tr>

<tr>
	<td colspan="5">
		<form class="grid" hx-post="/admin/clues/create" hx-confirm="unset"  hx-swap="outerHTML">
			<label for="location">
				Location
				<input type="text" id="location" name="location" placeholder="Location" required="">
			</label>
			
			<label for="clue">
				Clue
				<input type="text" id="clue" name="clue" placeholder="Clue" required="">
			</label>
			
			<button>Add</button>
		</form>
	</td>
</tr>
`,
		clue.Code, clue.Location, clue.Clue, clue.Code, clue.Code)
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

// ChangeClues shows a team all of the clues they have unlocked
func ChangeClues(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	code := chi.URLParam(r, "code")

	clue := models.Clue{}
	result := env.DB.Where("code = ?", code).Find(&clue)
	if result.Error != nil {
		return handler.StatusError{http.StatusNotFound, errors.New("this clue does not exist")}
	}

	if r.Method == http.MethodDelete {
		result = env.DB.Delete(&clue)
		if result.Error != nil {
			return handler.StatusError{http.StatusInternalServerError, errors.New("this clue could not be deleted")}
		} else {
			return nil
		}
	}

	return handler.StatusError{http.StatusMethodNotAllowed, errors.New("not yet implemented")}
}
