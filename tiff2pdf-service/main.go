package main

import (
	"crypto/rand"
	"encoding/base64"
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
	if newBind := os.Getenv("TIFF2PDF_SERVICE_LISTEN"); newBind != "" {
		bind = newBind
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

func failConversion(xReqId string, w http.ResponseWriter, err error) {
	log.Printf("[%s] %s", xReqId, err.Error())
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}

func convertTiff2PDF(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var b []byte
	var o *tiff2pdf.ConvertTiffToPDFOutput
	var err error

	xReqId := req.Header.Get("X-Request-ID")
	if len(xReqId) == 0 {
		rb := make([]byte, 6)
		_, err := rand.Read(rb)
		if err == nil {
			xReqId = base64.URLEncoding.EncodeToString(rb)
		} else {
			xReqId = "NONE"
		}
	}

	success := false

	defer func() {
		end := time.Now()
		diff := end.Sub(start)
		if success {
			log.Printf("[%s] Converted %d bytes TIFF to %d bytes PDF in %v", xReqId, len(b), len(o.PDF), diff)
		} else {
			log.Printf("[%s] Failed conversion of %d bytes TIFF to PDF in %v", xReqId, len(b), diff)
		}
	}()

	if b, err = ioutil.ReadAll(req.Body); err != nil {
		failConversion(xReqId, w, err)
		return
	}
	log.Printf("[%s] Got %d input TIFF bytes", xReqId, len(b))

	req.Body.Close()

	c := tiff2pdf.DefaultConfig()

	if hdr, ok := req.Header["PDF-PageSize"]; ok {
		log.Printf("[%s] Setting PDF page size: %s", xReqId, hdr[0])
		c.PageSize = hdr[0]
	}
	if hdr, ok := req.Header["PDF-FullPage"]; ok {
		log.Printf("[%s] Setting PDF full page: %s", xReqId, hdr[0])
		fullPage, err := strconv.ParseBool(hdr[0])
		if err != nil {
			failConversion(xReqId, w, err)
			return
		}
		c.FullPage = fullPage
	}
	if hdr, ok := req.Header["PDF-Subject"]; ok {
		log.Printf("[%s] Setting PDF subject: %s", xReqId, hdr[0])
		c.Subject = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Author"]; ok {
		log.Printf("[%s] Setting PDF author: %s", xReqId, hdr[0])
		c.Author = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Creator"]; ok {
		log.Printf("[%s] Setting PDF creator: %s", xReqId, hdr[0])
		c.Creator = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Title"]; ok {
		log.Printf("[%s] Setting PDF title: %s", xReqId, hdr[0])
		c.Title = hdr[0]
	}

	if o, err = tiff2pdf.ConvertTiffToPDF(b, c, "input.tif", "output.pdf"); err != nil {
		failConversion(xReqId, w, err)
		return
	}

	success = true

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("PDF-Pages", strconv.Itoa(int(o.PageCount)))
	w.WriteHeader(200)
	w.Write(o.PDF)
}
