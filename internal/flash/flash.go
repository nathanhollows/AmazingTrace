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
func Get(w http.ResponseWriter, r *http.Request) []Message {
	session, _ := getCookieStore().Get(r, sessionName)
	fm := session.Flashes("message")
	// If we have some messages.
	if len(fm) > 0 {
		session.Save(r, w)
		// Initiate a strings slice to return messages.
		var flashes []Message
		for _, fl := range fm {
			// Add message to the slice.
			flashes = append(flashes, fl.(Message))
		}

		return flashes
	}
	return nil
}
