package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/surakshith-suvarna/bookings/internal/config"
	"github.com/surakshith-suvarna/bookings/internal/models"
)

//This allows functions defined in this to be used by any templates
var functions = template.FuncMap{
	"formatDate":         FormatDate,
	"formatCalendarDate": FormatCalendarDate,
	"iterate":            Iterate,
	"add":                Add,
}

var app *config.AppConfig
var pathToTemplates = "./templates"

//NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

//Formats Date in string format YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

//FormatCalendarDate formats the date in format specified
func FormatCalendarDate(t time.Time, f string) string {
	return t.Format(f)
}

//Iterate can be used instead of for loops in template as Go templates do not have for loops
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

//Add adds two values within template
func Add(a, b int) int {
	return a + b
}

//AddDefaultData adds a data automatically on multiple pages (all Templates)
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

//Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.UseCache {
		//get the template cache from app config map instead of reading from disk
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		//log.Fatal("Could not get template from template cache")
		return errors.New("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	//adds default data to all templates
	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	//parsedTemplate, _ := template.ParseFiles("./templates/" + templ)
	//err = parsedTemplate.Execute(w, nil)
	return nil

}

//CreateTemplateCache creates a template cache
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	//pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println("page accessed is", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		//matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			//ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil

}
