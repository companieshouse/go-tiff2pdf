package main

import (
	"io/ioutil"
	"log"

	"github.com/companieshouse/go-tiff2pdf/tiff2pdf"
)

func main() {
	files, err := ioutil.ReadDir("tifs")
	if err != nil {
		log.Fatal(err)
	}
	fileCount, errorCount := 0, 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileCount++
		inputName := file.Name()
		b, err := ioutil.ReadFile("tifs/" + inputName)
		if err != nil {
			log.Fatalf("%s, reading %s after %d files with %d errors", err, inputName, fileCount, errorCount)
		}

		outputName := inputName + ".pdf"
		o, err := tiff2pdf.ConvertTiffToPDF(b, tiff2pdf.DefaultConfig(), inputName, outputName)
		if err != nil {
			errorCount++
			log.Printf("ERROR in %s: %s\n", inputName, err)
		}

		if err = ioutil.WriteFile("pdfs/"+outputName, o.PDF, 0644); err != nil {
			log.Fatalf("%s, writing %s after %d files with %d errors", err, outputName, fileCount, errorCount)
		}
	}
	log.Printf("Done %d files with %d errors\n", fileCount, errorCount)
}
