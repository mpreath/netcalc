package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Config *Config
}

func NewApp(config *Config) *App {
	app := &App{
		Router: mux.NewRouter(),
		Config: config,
	}

	app.initialize()

	return app
}

func (app *App) initialize() {
	app.Router.Path("/jwt").Methods(http.MethodGet).HandlerFunc(app.GetJWT)
	app.Router.Path("/info").Methods(http.MethodGet).HandlerFunc(app.ValidateJWT(Info))
	app.Router.Path("/subnet").Methods(http.MethodGet).HandlerFunc(app.ValidateJWT(Subnet))
	app.Router.Path("/summarize").Methods(http.MethodPost).HandlerFunc(app.ValidateJWT(Summarize))
	app.Router.Path("/vlsm").Methods(http.MethodGet).HandlerFunc(app.ValidateJWT(Vlsm))
}

func (app *App) Run() {
	log.Printf("Netcalc API Server Started [%s]\n", app.Config.HttpPort)
	log.Println(http.ListenAndServe(":"+app.Config.HttpPort, app.Router))
}
