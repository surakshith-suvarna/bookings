package main

import (
	"encoding/gob"
	"flag"
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
	//What type of data is stored in session (The data types allowed in session. If required add new ones here)
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	inProduction := flag.Bool("production", true, "App in Production")
	useCache := flag.Bool("cache", true, "Use cache for template")
	dbHost := flag.String("host", "localhost", "DB Host")
	dbName := flag.String("dbname", "", "DB Name")
	dbUser := flag.String("dbuser", "", "DB User")
	dbPass := flag.String("dbpass", "", "DB Pass")
	dbPort := flag.String("dbport", "5432", "DB Port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL Settings (disable, prefer, require)")

	flag.Parse()
	if *dbName == "" || *dbUser == "" || *dbPass == "" {
		fmt.Println("Required Flags are empty")
		os.Exit(1)
	}
	//Create channel for mailData
	mailChan := make(chan models.MailData)
	//PRovides access to channel across our app
	app.MailChan = mailChan

	//change this to true when in Production
	//app.InProduction = false
	app.InProduction = *inProduction

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
	connecttionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connecttionString)
	//db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=R00t@ss001")
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
	//app.UseCache = false
	app.UseCache = *useCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
