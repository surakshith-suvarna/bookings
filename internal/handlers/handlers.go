package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/surakshith-suvarna/bookings/internal/config"
	"github.com/surakshith-suvarna/bookings/internal/forms"
	"github.com/surakshith-suvarna/bookings/internal/helpers"
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

	//send data to the template
	render.RenderTemplates(w, r, "about.page.tmpl", &models.TemplateData{})
}

//Reservation renders the make a reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplates(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})

}

//PostReservation handles posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	//create form object
	forms := forms.New(r.PostForm)

	//forms.Has("first_name", r)
	forms.Required("first_name", "last_name", "email")
	forms.MinLength("first_name", 3)
	forms.IsEmail("email")

	if !forms.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplates(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: forms,
			Data: data,
		})

		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

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
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		//Sets Error data to session and redirects to homepage
		m.App.Session.Put(r.Context(), "error", "Reservation details not found in session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//Delete the reservation data from session
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplates(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
