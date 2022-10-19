package main

import (
	"github.com/gorilla/mux"
)

func main() {

	app := &App{}

	app.Router = mux.NewRouter()
	app.InitializeRoutes()

	app.Run()
}
