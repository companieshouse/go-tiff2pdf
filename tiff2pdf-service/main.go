package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/companieshouse/go-tiff2pdf/tiff2pdf"
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
	p.Path("/convert").Methods("POST").HandlerFunc(convertTiff2PDF)

	http.Handle("/", p)
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func convertTiff2PDF(w http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	req.Body.Close()

	c := tiff2pdf.DefaultConfig()
	o, err := tiff2pdf.ConvertTiff2PDF(b, c, "input.tif", "output.pdf")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(o)
}
