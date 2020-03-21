package web

import (
	"github.com/markbates/pkger"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	TemplatesDir = pkger.Dir("/templates")
	templates = make(map[string]*template.Template)
)

func ExecuteTemplate(name, path string, f func(r *http.Request) interface {}) http.HandlerFunc {
	temp, ok := templates[name]
	if !ok {
		file, err := TemplatesDir.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		text, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		t, err := template.New(name).Parse(string(text))
		templates[name] = t
		temp = t
	}
	return func(w http.ResponseWriter, r * http.Request) {
		err := temp.Execute(w, f(r))
		if err != nil {
			log.Fatal(err)
		}
	}
}

type IndexVars struct {}
func Index(r *http.Request) interface{} {
	return IndexVars{}
}
