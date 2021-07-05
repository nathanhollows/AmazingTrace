package flash

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

const sessionName = "flashMessages"

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
	// TODO: In real-world applications, use env variables to store the session key.
	sessionKey := "test-session-key"
	return sessions.NewCookieStore([]byte(sessionKey))
}

// Set adds a new message into the cookie storage.
func Set(w http.ResponseWriter, r *http.Request, message Message) {
	session, _ := getCookieStore().Get(r, sessionName)
	session.AddFlash(message, "message")
	session.Save(r, w)
}

// Get gets flash messages from the cookie storage.
func Get(session *sessions.Session, w http.ResponseWriter, r *http.Request) []interface{} {
	messages := session.Flashes()
	session.Save(r, w)
	return messages
}
