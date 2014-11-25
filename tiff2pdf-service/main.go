package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/companieshouse/go-tiff2pdf/tiff2pdf"
	"github.com/gorilla/pat"
)

func main() {
	setupRouting()

	bind := ":9090"
	if os.Getenv("TIFF2PDF_SERVICE_LISTEN") != "" {
		bind = os.Getenv("TIFF2PDF_SERVICE_LISTEN")
	}
	log.Printf("Listening on %s", bind)

	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal(err)
	}
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

func failConversion(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}

func convertTiff2PDF(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var b, o []byte
	var err error

	defer func() {
		end := time.Now()
		diff := end.Sub(start)
		log.Printf("Converted %d bytes TIFF to %d bytes PDF in %v", len(b), len(o), diff)
	}()

	if b, err = ioutil.ReadAll(req.Body); err != nil {
		failConversion(w, err)
		return
	}

	req.Body.Close()

	c := tiff2pdf.DefaultConfig()

	if hdr, ok := req.Header["PDF-PageSize"]; ok {
		log.Printf("Setting PDF page size: %s", hdr[0])
		c.PageSize = hdr[0]
	}
	if hdr, ok := req.Header["PDF-FullPage"]; ok {
		log.Printf("Setting PDF full page: %s", hdr[0])
		fullPage, err := strconv.ParseBool(hdr[0])
		if err != nil {
			failConversion(w, err)
			return
		}
		c.FullPage = fullPage
	}
	if hdr, ok := req.Header["PDF-Subject"]; ok {
		log.Printf("Setting PDF subject: %s", hdr[0])
		c.Subject = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Author"]; ok {
		log.Printf("Setting PDF author: %s", hdr[0])
		c.Author = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Creator"]; ok {
		log.Printf("Setting PDF creator: %s", hdr[0])
		c.Creator = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Title"]; ok {
		log.Printf("Setting PDF title: %s", hdr[0])
		c.Title = hdr[0]
	}

	if o, err = tiff2pdf.ConvertTiffToPDF(b, c, "input.tif", "output.pdf"); err != nil {
		failConversion(w, err)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(o)
}
