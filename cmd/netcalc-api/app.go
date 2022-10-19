package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

func (app *App) InitializeRoutes() {
	app.Router.Path("/info").Methods(http.MethodGet).HandlerFunc(Info)
	app.Router.Path("/subnet").Methods(http.MethodGet).HandlerFunc(Subnet)
	app.Router.Path("/summarize").Methods(http.MethodPost).HandlerFunc(Summarize)
	app.Router.Path("/vlsm").Methods(http.MethodGet).HandlerFunc(Vlsm)
}

func (app *App) Run() {
	log.Println("Start listening")
	log.Println(http.ListenAndServe(":8080", app.Router))
}
