package helpers

import (
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// URL constructs a URL specific to the application
func URL(patterns ...string) string {
	u := &url.URL{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	site := os.Getenv("TRACE_URL")
	if site != "" {
		u, _ = url.Parse(site)
	} else {
		u.Path = "/"
	}
	if len(patterns) > 0 {
		u.Path += patterns[0]
	}
	if len(patterns) > 1 {
		u.RawQuery = patterns[1]
	}
	if len(patterns) > 2 {
		u.Fragment = patterns[2]
	}
	return u.String()
}
