package main

import (
	"fmt"
	"io/ioutil"
	"log"

	tiff2pdf "github.com/companieshouse/go-tiff2pdf"
)

func main() {
	b, err := ioutil.ReadFile("test.tif")
	if err != nil {
		log.Fatal(err)
	}

	bOut, err := tiff2pdf.ConvertTiffToPDF(b)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("success?")
	fmt.Printf("%s", bOut)
}
