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

func failConversion(xReqId string, w http.ResponseWriter, err error) {
	log.Printf("%s%s", xReqId, err.Error())
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}

func convertTiff2PDF(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	var b []byte
	var o *tiff2pdf.ConvertTiffToPDFOutput
	var err error

	xReqId := req.Header.Get("X-Request-ID")
	if len(xReqId) > 0 {
		xReqId = "[" + xReqId + "] "
	}

	success := false

	defer func() {
		end := time.Now()
		diff := end.Sub(start)
		if success {
			log.Printf("%sConverted %d bytes TIFF to %d bytes PDF in %v", xReqId, len(b), len(o.PDF), diff)
		} else {
			log.Printf("%sFailed conversion of %d bytes TIFF to PDF in %v", xReqId, len(b), diff)
		}
	}()

	if b, err = ioutil.ReadAll(req.Body); err != nil {
		failConversion(xReqId, w, err)
		return
	}
	log.Printf("%sGot %d input TIFF bytes", xReqId, len(b))

	req.Body.Close()

	c := tiff2pdf.DefaultConfig()

	if hdr, ok := req.Header["PDF-PageSize"]; ok {
		log.Printf("%sSetting PDF page size: %s", xReqId, hdr[0])
		c.PageSize = hdr[0]
	}
	if hdr, ok := req.Header["PDF-FullPage"]; ok {
		log.Printf("%sSetting PDF full page: %s", xReqId, hdr[0])
		fullPage, err := strconv.ParseBool(hdr[0])
		if err != nil {
			failConversion(xReqId, w, err)
			return
		}
		c.FullPage = fullPage
	}
	if hdr, ok := req.Header["PDF-Subject"]; ok {
		log.Printf("%sSetting PDF subject: %s", xReqId, hdr[0])
		c.Subject = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Author"]; ok {
		log.Printf("%sSetting PDF author: %s", xReqId, hdr[0])
		c.Author = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Creator"]; ok {
		log.Printf("%sSetting PDF creator: %s", xReqId, hdr[0])
		c.Creator = hdr[0]
	}
	if hdr, ok := req.Header["PDF-Title"]; ok {
		log.Printf("%sSetting PDF title: %s", xReqId, hdr[0])
		c.Title = hdr[0]
	}

	if o, err = tiff2pdf.ConvertTiffToPDF(b, c, "input.tif", "output.pdf"); err != nil {
		failConversion(xReqId, w, err)
		return
	}

	success = true

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("PDF-Pages", strconv.Itoa(int(o.PageCount)))
	w.Write(o.PDF)
}
