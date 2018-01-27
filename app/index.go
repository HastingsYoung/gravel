package web

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var templates = template.Must(template.ParseFiles(makePath("index.html")))
var validPath = regexp.MustCompile("^/(index)*$")

func Index(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "index")
	return
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makePath(path string) string {
	return filepath.Join(os.Getenv("GOPATH")+"/src/github.com/gravel/"+"app", "template", path)
}
