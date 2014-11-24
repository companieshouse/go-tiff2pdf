package main

import (
	//"fmt"
	"io/ioutil"
	"log"

	tiff2pdf "github.com/companieshouse/go-tiff2pdf"
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
			log.Fatal(err)
		}

		outputName := inputName + ".pdf"
		bOut, err := tiff2pdf.ConvertTiffToPDF(b, tiff2pdf.DefaultConfig(), inputName, outputName)
		if err != nil {
			errorCount++
			log.Printf("ERROR in %s: %s\n", inputName, err)
		}

		// fmt.Printf("%s", bOut)
		if err = ioutil.WriteFile("pdfs/"+outputName, bOut, 0644); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Done %d files with %d errors\n", fileCount, errorCount)
}
