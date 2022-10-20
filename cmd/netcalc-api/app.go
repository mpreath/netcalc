package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Config *Configuration
}

func NewApp(config *Configuration) *App {
	app := &App{
		Router: mux.NewRouter(),
		Config: config,
	}

	app.initialize()

	return app
}

func (app *App) initialize() {
	app.Router.Path("/info").Methods(http.MethodGet).HandlerFunc(Info)
	app.Router.Path("/subnet").Methods(http.MethodGet).HandlerFunc(Subnet)
	app.Router.Path("/summarize").Methods(http.MethodPost).HandlerFunc(Summarize)
	app.Router.Path("/vlsm").Methods(http.MethodGet).HandlerFunc(Vlsm)
}

func (app *App) Run() {
	log.Printf("Netcalc API Server Started [%s]\n", strconv.Itoa(app.Config.HttpPort))
	log.Println(http.ListenAndServe(":"+strconv.Itoa(app.Config.HttpPort), app.Router))
}
