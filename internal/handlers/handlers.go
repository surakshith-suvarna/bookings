package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/surakshith-suvarna/bookings/internal/config"
	"github.com/surakshith-suvarna/bookings/internal/models"
	"github.com/surakshith-suvarna/bookings/internal/render"
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
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "home.page.tmpl", &models.TemplateData{})
}

//About is the about page holder
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic

	stringMapValue := make(map[string]string)

	stringMapValue["test"] = "This is test data"
	//send data to the template
	render.RenderTemplates(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMapValue,
	})
}

//Reservation renders the make a reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

//Generals renders the generals-quaters page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "generals.page.tmpl", &models.TemplateData{})
}

//Majors renders the majors-suite page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "majors.page.tmpl", &models.TemplateData{})
}

//Contact renders the Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "contact.page.tmpl", &models.TemplateData{})
}

//Availability renders the search-availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

//PostAvailability processes the search-availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	//render.RenderTemplates(w, "search-availability.page.tmpl", &models.TemplateData{})
	start := r.Form.Get("start")
	end := r.Form.Get(("end"))

	w.Write([]byte(fmt.Sprintf("The value of start and end dates is %s and %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

//AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	//render.RenderTemplates(w, "search-availability.page.tmpl", &models.TemplateData{})
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
