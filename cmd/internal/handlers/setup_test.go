package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	scs "github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
	"github.com/maslow123/bookings/cmd/internal/config"
	"github.com/maslow123/bookings/cmd/internal/models"
	"github.com/maslow123/bookings/cmd/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../../templates"
var functions = template.FuncMap{}
var infoLog *log.Logger
var errorLog *log.Logger

func getRoutes() http.Handler {
	// what am i going to put in the session
	gob.Register(models.Reservation{})
	// change this true when in production

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

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := NewRepo(&app)

	NewHandlers(repo)

	render.NewTemplates(&app)

	mux := chi.NewRouter()

	// mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", http.HandlerFunc(Repo.Home))
	mux.Get("/about", http.HandlerFunc(Repo.About))
	mux.Get("/generals-quarters", http.HandlerFunc(Repo.Generals))
	mux.Get("/majors-suite", http.HandlerFunc(Repo.Majors))
	mux.Get("/search-availability", http.HandlerFunc(Repo.Availability))
	mux.Post("/search-availability", http.HandlerFunc(Repo.PostAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(Repo.AvailabilityJSON))
	mux.Get("/make-reservation", http.HandlerFunc(Repo.Reservation))
	mux.Post("/make-reservation", http.HandlerFunc(Repo.PostReservation))
	mux.Get("/contact", http.HandlerFunc(Repo.Contact))
	mux.Get("/reservation-summary", http.HandlerFunc(Repo.ReservationSummary))

	fileServer := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/assets/*", http.StripPrefix("/assets", fileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.htm", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.htm", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.htm", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
