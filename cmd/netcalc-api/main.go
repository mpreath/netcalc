package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	r.Path("/info").Methods(http.MethodGet).HandlerFunc(Info)
	r.Path("/subnet").Methods(http.MethodGet).HandlerFunc(Subnet)
	r.Path("/summarize").Methods(http.MethodPost).HandlerFunc(Summarize)
	r.Path("/vlsm").Methods(http.MethodGet).HandlerFunc(Vlsm)
	log.Println("Start listening")
	log.Println(http.ListenAndServe(":8080", r))
}
