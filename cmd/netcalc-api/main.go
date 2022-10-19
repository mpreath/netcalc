package main

import (
	"github.com/gorilla/mux"
)

func main() {

	app := NewApp()
	app.Router = mux.NewRouter()
	app.InitializeRoutes()
	app.Run()
}
