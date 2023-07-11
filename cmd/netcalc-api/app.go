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
	// Middleware
	app.Router.Use(LoggingMiddleware)
	app.Router.Use(CORSMiddleware)

	// Routes
	jwtRouter := app.Router.PathPrefix("/token").Subrouter()
	jwtRouter.Path("/new").Methods(http.MethodGet).HandlerFunc(app.GetJWT)

	apiRouter := app.Router.PathPrefix("/api").Subrouter()
	// apiRouter.Use(app.ValidateJWT)
	apiRouter.Path("/info").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(Info)
	apiRouter.Path("/subnet").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(Subnet)
	apiRouter.Path("/summarize").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(Summarize)
	apiRouter.Path("/vlsm").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(Vlsm)
}

func (app *App) Run() {
	log.Printf("Netcalc API Server Started [%s]\n", app.Config.HttpPort)
	log.Println(http.ListenAndServe(":"+app.Config.HttpPort, app.Router))
}
