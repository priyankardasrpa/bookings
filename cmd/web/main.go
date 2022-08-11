package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/priyankardasrpa/bookings/pkg/config"
	"github.com/priyankardasrpa/bookings/pkg/handlers"
	"github.com/priyankardasrpa/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

var portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Server is running on port %s\n", portNumber)
	//http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
