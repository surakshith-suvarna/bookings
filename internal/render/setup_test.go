package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/surakshith-suvarna/bookings/internal/config"
	"github.com/surakshith-suvarna/bookings/internal/models"
)

var testApp config.AppConfig
var session *scs.SessionManager

func TestMain(m *testing.M) {

	//What type of data is stored in session
	gob.Register(models.Reservation{})

	//change this to true when in Production
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())

}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil

}

func (tw *myWriter) WriteHeader(i int) {

}
