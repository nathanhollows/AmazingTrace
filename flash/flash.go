package flash

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(Message{})
}

// Message is a struct containing each flashed message
type Message struct {
	Title   string
	Message string
	Style   string
}

func getCookieStore() *sessions.CookieStore {
	key := []byte(os.Getenv("SESSION_KEY"))
	return sessions.NewCookieStore(key)
}

// Set adds a new message into the cookie storage.
func Set(w http.ResponseWriter, r *http.Request, message Message) {
	store := getCookieStore()
	session, _ := store.Get(r, "flash")
	session.Options.HttpOnly = true
	session.Options.SameSite = http.SameSiteStrictMode
	// session.Options.Secure = true
	session.AddFlash(message)
	session.Save(r, w)
}

// Success adds a new success message into the cookie storage.
func Success(w http.ResponseWriter, r *http.Request, message string) {
	Set(w, r, Message{Title: "Success", Message: message, Style: "success"})
}

// Error adds a new error message into the cookie storage.
func Error(w http.ResponseWriter, r *http.Request, message string) {
	Set(w, r, Message{Message: message, Style: "danger"})
}

// Warning adds a new warning message into the cookie storage.
func Warning(w http.ResponseWriter, r *http.Request, message string) {
	Set(w, r, Message{Title: "Warning", Message: message, Style: "warning"})
}

// Info adds a new info message into the cookie storage.
func Info(w http.ResponseWriter, r *http.Request, message string) {
	Set(w, r, Message{Title: "Info", Message: message, Style: "info"})
}

// Get flash messages from the cookie storage.
func Get(w http.ResponseWriter, r *http.Request) []interface{} {
	session, err := getCookieStore().Get(r, "flash")
	if err == nil {
		messages := session.Flashes()
		if len(messages) > 0 {
			session.Options.MaxAge = -1
			session.Save(r, w)
		}
		return messages
	}
	return nil
}
