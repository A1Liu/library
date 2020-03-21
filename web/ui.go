package web

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"html/template"
	"io/ioutil"
	"log"
)

var (
	TemplatesDir = pkger.Dir("/templates")
	templates = make(map[string]*template.Template)
)

func ExecuteTemplate(name, path string, f func(c *gin.Context) interface {}) gin.HandlerFunc {
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
	return func(c *gin.Context) {
		err := temp.Execute(c.Writer, f(c))
		if err != nil {
			log.Fatal(err)
		}
	}
}

type IndexVars struct {}
func Index(c *gin.Context) interface{} {
	return IndexVars{}
}
