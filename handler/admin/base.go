package admin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
	"divide": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b)
	},
	"progress": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b) * 100
	},
	"add": func(a, b int) int {
		return a + b
	},
	"date": func() time.Time {
		return time.Now()
	},
}

func parse(patterns ...string) *template.Template {
	patterns = append(patterns, "layout.html", "flash.html")
	for i := 0; i < len(patterns); i++ {
		patterns[i] = "templates/admin/" + patterns[i]
	}
	return template.Must(template.New("base").Funcs(funcs).ParseFiles(patterns...))
}

func render(w http.ResponseWriter, data map[string]interface{}, patterns ...string) error {
	w.Header().Set("Content-Type", "text/html")
	if data["title"] != nil {
		data["title"] = fmt.Sprintf("%v | Admin", data["title"])
	} else {
		data["title"] = "Admin"
	}
	err := parse(patterns...).ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return err
}

func renderFragment(w http.ResponseWriter, data map[string]interface{}, patterns ...string) error {
	w.Header().Set("Content-Type", "text/html")

	patterns = append(patterns, "flash.html")
	for i := 0; i < len(patterns); i++ {
		patterns[i] = "templates/admin/" + patterns[i]
	}
	err := template.Must(template.New("base").Funcs(funcs).ParseFiles(patterns...)).ExecuteTemplate(w, "fragment", data)

	if err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return err
}
