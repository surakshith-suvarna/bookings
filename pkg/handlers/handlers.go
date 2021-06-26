package handlers

import (
	"net/http"

	"github.com/surakshith-suvarna/bookings/pkg/config"
	"github.com/surakshith-suvarna/bookings/pkg/models"
	"github.com/surakshith-suvarna/bookings/pkg/render"
)

//Repo the repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page holder
func (m *Repository) Home(w http.ResponseWriter, h *http.Request) {
	render.RenderTemplates(w, "home.page.tmpl", &models.TemplateData{})
}

//About is the about page holder
func (m *Repository) About(w http.ResponseWriter, h *http.Request) {
	//perform some logic

	stringMapValue := make(map[string]string)

	stringMapValue["test"] = "This is test data"
	//send data to the template
	render.RenderTemplates(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMapValue,
	})
}
