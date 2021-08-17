package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	scs "github.com/alexedwards/scs/v2"
	"github.com/maslow123/bookings/pkg/config"
	"github.com/maslow123/bookings/pkg/handlers"
	"github.com/maslow123/bookings/pkg/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	// change this true when in production

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println("Starting application on port", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}