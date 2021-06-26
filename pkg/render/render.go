package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/surakshith-suvarna/bookings/pkg/config"
	"github.com/surakshith-suvarna/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

//NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

//AddDefaultData adds a data automatically on multiple pages
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

//RenderTemplate renders templates using html/template
func RenderTemplates(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		//get the template cache from app config map instead of reading from disk
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	//adds default data
	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

	//parsedTemplate, _ := template.ParseFiles("./templates/" + templ)
	//err = parsedTemplate.Execute(w, nil)

}

//CreateTemplateCache creates a template cache
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		fmt.Println("page accessed is", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil

}
