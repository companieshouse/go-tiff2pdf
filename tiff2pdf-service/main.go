package main

import (
	"log"
	"net/http"

	"github.com/gorilla/pat"
)

func main() {
	setupRouting()
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan int)
}

func setupRouting() {
	p := pat.New()

	p.Path("/healthcheck").Methods("GET").HandlerFunc(healthcheck)

	http.Handle("/", p)
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}
