package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/AmazingTrace/flash"
	"github.com/nathanhollows/AmazingTrace/handler"
	"github.com/nathanhollows/AmazingTrace/models"
)

// Clues lists all available clues
func Clues(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	clues := []models.Clue{}
	result := env.DB.Find(&clues)

	if result.RowsAffected > 0 {
		env.Data["clues"] = clues
	}
	env.Data["title"] = "Clues"
	env.Data["messages"] = flash.Get(w, r)
	return render(w, env.Data, "clues/index.html")
}

// NewClue shows a form to create a new clue
func NewClue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	env.Data["title"] = "New Clue"

	if r.Method == http.MethodPost {
		r.ParseForm()
		clue := models.Clue{}
		clue.Location = r.PostFormValue("name")
		clue.Clue = r.PostFormValue("clues")
		clue.Longitude = r.PostFormValue("longitude")
		clue.Latitude = r.PostFormValue("latitude")
		points, _ := strconv.Atoi(r.PostFormValue("points"))
		clue.Points = points
		result := env.DB.Create(&clue)
		if result.Error != nil {
			flash.Set(w, r, flash.Message{Message: "Could not save clue", Style: "warning"})
		} else {
			message := fmt.Sprintf("Clue <a href=\"/admin/clues/edit/%v\">%v</a> for %v successfully saved", clue.Code, clue.Code, clue.Location)
			flash.Set(w, r, flash.Message{Message: message, Style: "success"})
		}
	}

	env.Data["messages"] = flash.Get(w, r)
	return render(w, env.Data, "clues/new.html")
}

// EditClue shows the edit form
func EditClue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	clueID := chi.URLParam(r, "code")
	clue := models.Clue{}
	result := env.DB.Where("code = ?", clueID).Find(&clue)
	if result.Error != nil {
		http.Redirect(w, r, "/admin/clues", http.StatusSeeOther)
		return handler.StatusError{Code: http.StatusNotFound, Err: errors.New("this clue does not exist")}
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		clue.Location = r.PostFormValue("name")
		clue.Clue = r.PostFormValue("clues")
		clue.Longitude = r.PostFormValue("longitude")
		clue.Latitude = r.PostFormValue("latitude")
		points, _ := strconv.Atoi(r.PostFormValue("points"))
		clue.Points = points
		result := env.DB.Save(&clue)
		if result.Error != nil {
			flash.Set(w, r, flash.Message{Message: "Could not save clue", Style: "warning"})
		} else {
			flash.Set(w, r, flash.Message{Message: "Clue successfully saved", Style: "success"})
		}
	}

	env.Data["title"] = "Editing " + clue.Location
	env.Data["clue"] = clue
	env.Data["messages"] = flash.Get(w, r)
	return render(w, env.Data, "clues/edit.html")
}

// CreateClue saves the posted clue
func CreateClue(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")

	if r.Method != http.MethodPost {
		return handler.StatusError{Code: http.StatusMethodNotAllowed, Err: errors.New("method not allowed")}
	}

	r.ParseForm()
	clue := models.Clue{}
	clue.Location = r.PostFormValue("location")
	clue.Clue = r.PostFormValue("clue")

	result := env.DB.Create(&clue)

	clue.GeneratePoster()

	if result.Error != nil {
		return handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("could not save clue")}
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
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}

	r.ParseForm()
	clue := models.Clue{}
	result := env.DB.Where("code = ?", r.PostFormValue("code")).Find(&clue)
	if result.Error != nil {
		session.AddFlash(flash.Message{Message: "Could not delete clue", Style: "warning"})
		session.Save(r, w)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}

	result = env.DB.Where("code = ?", r.PostFormValue("code")).Delete(&clue)
	if result.Error != nil {
		session.AddFlash(flash.Message{Message: "Could not delete clue", Style: "warning"})
		session.Save(r, w)
	} else {
		session.AddFlash(flash.Message{Message: "Clue successfully deleted", Style: "success"})
		session.Save(r, w)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	return nil
}

// ChangeClues shows a team all of the clues they have unlocked
func ChangeClues(env *handler.Env, w http.ResponseWriter, r *http.Request) error {
	code := chi.URLParam(r, "code")

	clue := models.Clue{}
	result := env.DB.Where("code = ?", code).Find(&clue)
	if result.Error != nil {
		return handler.StatusError{Code: http.StatusNotFound, Err: errors.New("this clue does not exist")}
	}

	if r.Method == http.MethodDelete {
		result = env.DB.Delete(&clue)
		if result.Error != nil {
			return handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("this clue could not be deleted")}
		} else {
			return nil
		}
	}

	return handler.StatusError{Code: http.StatusMethodNotAllowed, Err: errors.New("not yet implemented")}
}
