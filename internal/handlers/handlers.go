package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/surakshith-suvarna/bookings/internal/config"
	"github.com/surakshith-suvarna/bookings/internal/driver"
	"github.com/surakshith-suvarna/bookings/internal/forms"
	"github.com/surakshith-suvarna/bookings/internal/helpers"
	"github.com/surakshith-suvarna/bookings/internal/models"
	"github.com/surakshith-suvarna/bookings/internal/render"
	"github.com/surakshith-suvarna/bookings/internal/repository"
	"github.com/surakshith-suvarna/bookings/internal/repository/dbrepo"
)

//Repo the repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//NewRepo creates a new Repository (This is available for all handlers appconfig and db pool)
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

//NewTestRepo creates a new Test Repository (This is only for tests)
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page holder
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

//About is the about page holder
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic

	//send data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

//Reservation renders the make a reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	//var emptyReservation models.Reservation
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//helpers.ServerError(w, errors.New("type assertion failed"))
		//return
		m.App.Session.Put(r.Context(), "error", "unable to get reservation data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//Getting RoomName
	room, err := m.DB.GetRoomById(res.RoomId)
	if err != nil {
		//helpers.ServerError(w, err)
		//return
		m.App.Session.Put(r.Context(), "error", "unable to fetch room data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName

	//Place res data along with RoomName in session for use in reservation summary page
	m.App.Session.Put(r.Context(), "reservation", res)

	//Since StartDate and EndDate are in time.Time format we cannot use in template. Hence convert to string
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})

}

//PostReservation handles posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	//Get session data for reservation which already contains start_date, end_date, roomId and roomName
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//helpers.ServerError(w, errors.New("type assertion failed"))
		//return
		m.App.Session.Put(r.Context(), "error", errors.New("type assertion failed"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	err := r.ParseForm()
	if err != nil {
		//helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "cant parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//sd := r.Form.Get("start_date")
	//ed := r.Form.Get("end_date")

	//2021-01-01-- 01/02 03:04:05PM '06 -0700 (GO Date format in Us Mountain Time)
	//layout := "2006-01-02"
	//startDate, err := time.Parse(layout, sd)
	//if err != nil {
	//		helpers.ServerError(w, err)
	//		return
	//	}
	//	endDate, err := time.Parse(layout, ed)
	//	if err != nil {
	//		helpers.ServerError(w, err)
	//		return
	//	}

	//	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	//	if err != nil {
	//		helpers.ServerError(w, err)
	//		return
	//	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	//	reservation := models.Reservation{
	//		FirstName: r.Form.Get("first_name"),
	//		LastName:  r.Form.Get("last_name"),
	//		Email:     r.Form.Get("email"),
	//		Phone:     r.Form.Get("phone"),
	//		StartDate: startDate,
	//		EndDate:   endDate,
	//		RoomId:    roomId,
	//	}

	//create form object
	forms := forms.New(r.PostForm)

	//forms.Has("first_name", r)
	forms.Required("first_name", "last_name", "email")
	forms.MinLength("first_name", 3)
	forms.IsEmail("email")

	if !forms.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: forms,
			Data: data,
		})

		return
	}

	//If data is valid insert the data into the database
	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		//helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		//helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//Send email to user
	htmlMessage := fmt.Sprintf(`<strong>Reservation Confirmed</strong></br>
					</br>
					<p>Dear %s</p></br>
					<p>Kindly find your reservation details below</br>
					<ul>
						<li>Room Type: %s</li>
						<li>Arrival: %s</li>
						<li>Departure: %s</li>
					</ul></p>`, reservation.FirstName, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"),
		reservation.EndDate.Format("2006-01-02"))
	userMail := models.MailData{
		From:     "bookings@domain.com",
		To:       reservation.Email,
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}
	m.App.MailChan <- userMail

	//Sent email to property owner
	ownerMsg := fmt.Sprintf(`
				<strong>Reservation  Confirmed</strong></br>
				<p>Room '%s' has been booked from %s to %s by %s</p>`, reservation.Room.RoomName,
		reservation.StartDate, reservation.EndDate, reservation.FirstName)

	ownerMail := models.MailData{
		From:    "owner@domain.com",
		To:      "owner@domain.com",
		Subject: "Reservation Booked",
		Content: ownerMsg,
	}
	m.App.MailChan <- ownerMail

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

//Generals renders the generals-quaters page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

//Majors renders the majors-suite page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

//Contact renders the Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

//Availability renders the search-availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

//PostAvailability processes the search-availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get(("end"))

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		// No availability
		m.App.Session.Put(r.Context(), "error", "No rooms available")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-page.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomId    string `json:"room_id"`
}

//AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	//render.RenderTemplates(w, "search-availability.page.tmpl", &models.TemplateData{})

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    strconv.Itoa(roomId),
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//ReservationSummary displays the reservation summary page
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

	//Date is currently in time.Time format and has to be converted to string to be used on Reservation Summary page
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

//ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	res.RoomId = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

//BookRoom takes url parameters, builds a sessional variable and takes user to make reservation page
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	//Get the data from GET request
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room, err := m.DB.GetRoomById(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomID,
	}

	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}
