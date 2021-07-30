package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/surakshith-suvarna/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quaters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	/*{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2021-07-13"},
		{key: "end", value: "2021-07-14"},
	}, http.StatusOK},
	{"post-avail-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2021-07-13"},
		{key: "end", value: "2021-07-14"},
	}, http.StatusOK},
	{"make reservation post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Surakshith"},
		{key: "last_name", value: "Suvarna"},
		{key: "email", value: "s@test.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},*/
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		//If get method else post method
		//if e.method == "GET" {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
		/*} else {
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}*/
	}
}

func TestRepository_Reservation(t *testing.T) {
	//Reservation function gets session in models.Reservation and hence create it
	reservation := models.Reservation{
		RoomId: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quaters",
		},
	}

	//Create a Request
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	//getting the context
	ctx := getCtx(req)
	//adding context to request
	req = req.WithContext(ctx)

	//Use a response recorder to mimic the response process. This is built into Go test package
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	//Take handler and convert to handler function. This will allow taking Reservation handler function to call directly
	handler := http.HandlerFunc(Repo.Reservation)

	//This will serve HTTP without calling routes as we have built it manually
	handler.ServeHTTP(rr, req)

	//Test for error
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code. Expected %d but received %d", http.StatusOK, rr.Code)
	}

	//Testing when reservation data is missing in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code. Expected %d but received %d", http.StatusTemporaryRedirect, rr.Code)
	}

	//Testing when RoomID is invalid
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomId = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code. Expected %d but received %d", http.StatusTemporaryRedirect, rr.Code)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	//All data are valid. Demostrates manual construction of post data
	//reqBody := "start_date=2050-01-02"
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-03")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=john")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=wick")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "email=wick@john.com")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-555-1234")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	//ctx := getCtx(req)
	//req = req.WithContext(ctx)
	//Tells this request this is a form request
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//rr := httptest.NewRecorder()
	//handler := http.HandlerFunc(Repo.PostReservation)
	//handler.ServeHTTP(rr, req)

	//if rr.Code != http.StatusTemporaryRedirect {
	//	t.Errorf("reservation handler returned wrong response code. Expected %d but received %d", http.StatusTemporaryRedirect, rr.Code)
	//}

	//test for invalid data
	/*layout := "2006-01-02"
	sd, _ := time.Parse(layout, "2050-01-01")
	ed, _ := time.Parse(layout, "2050-01-02")

	reservation := models.Reservation{
		StartDate: sd,
		EndDate:   ed,
		RoomId:    1,
	}

	reqBody := fmt.Sprintf("start_date=%s", reservation.StartDate.Format("2006-01-02"))
	reqBody = fmt.Sprintf("%s&end_date=%s", reqBody, reservation.EndDate.Format("2006-01-02"))
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=j")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=wick")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=wick")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-555-1234")
	reqBody = fmt.Sprintf("%s&room_id=%d", reqBody, reservation.RoomId)*/

	//Testing session get
	layout := "2006-01-02"
	sd, _ := time.Parse(layout, "2050-01-01")
	ed, _ := time.Parse(layout, "2050-01-02")
	reservation := models.Reservation{
		StartDate: sd,
		FirstName: "John",
		LastName:  "Wick",
		Email:     "John@wick.com",
		Phone:     "444-555-1234",
		EndDate:   ed,
		RoomId:    1,
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(req.Context(), "reservation", reservation)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returned wrong response code. Expected %d but received %d", http.StatusOK, rr.Code)
	}

	//Testing when no data is posted
	//reservebyte, _ := json.Marshal(reservation)
	//reader := bytes.NewReader(reservebyte)
	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "wick")
	postedData.Add("email", "John@wick.com")
	postedData.Add("phone", "555-555-1234")
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-post-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code. Expected %d but received %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// Test for Insert restrictions table
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "wick")
	postedData.Add("email", "John@wick.com")
	postedData.Add("phone", "555-555-1234")
	postedData.Add("room_id", "2")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-post-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code. Expected %d but received %d", http.StatusTemporaryRedirect, rr.Code)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
