package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cyberjourney20/bookings/internal/config"
	"github.com/cyberjourney20/bookings/internal/driver"
	"github.com/cyberjourney20/bookings/internal/handlers"
	"github.com/cyberjourney20/bookings/internal/helpers"
	"github.com/cyberjourney20/bookings/internal/models"
	"github.com/cyberjourney20/bookings/internal/render"
	"github.com/joho/godotenv"

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
	//Close db connectiion when program closes
	defer db.SQL.Close()
	defer close(app.MailChan)

	fmt.Println("Staring mail listener...")
	listenForMail()

	fmt.Printf("Starting application on port: %s\n", portNumber)

	// Starts HTTP Server
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// What goes in the session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//change to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//godotenv to load secrets from .env file
	err := godotenv.Load(os.ExpandEnv("./.env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbString := os.Getenv("DBSTRING")

	//connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(dbString)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect to database! Dying... %v\n", err))
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
