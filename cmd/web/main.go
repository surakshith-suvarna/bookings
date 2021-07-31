package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/surakshith-suvarna/bookings/internal/config"
	"github.com/surakshith-suvarna/bookings/internal/driver"
	"github.com/surakshith-suvarna/bookings/internal/handlers"
	"github.com/surakshith-suvarna/bookings/internal/helpers"
	"github.com/surakshith-suvarna/bookings/internal/models"
	"github.com/surakshith-suvarna/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	//Defer Mail Channel to close it when program stops executing
	defer close(app.MailChan)

	//Listen to the mail channel
	fmt.Println("Starting Mail Listener")
	listenForMail()

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	//http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	//What type of data is stored in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	//Create channel for mailData
	mailChan := make(chan models.MailData)
	//PRovides access to channel across our app
	app.MailChan = mailChan

	//change this to true when in Production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//Initialize DB Repo (Connect to DB)
	log.Println("connecting to datbase....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=R00t@ss001")
	if err != nil {
		log.Fatal("Cannot connect to database..")
	}
	log.Println("Connected to Database...")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
