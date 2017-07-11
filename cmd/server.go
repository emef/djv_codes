package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/emef/djv_codes"
)

func main() {
	port := flag.String("port", ":8080", "Port http service will run on")
	codesDir := flag.String(
		"codes_dir", "codes", "Directory where all available codes files live")
	usedCodesPath := flag.String(
		"used_codes_file", "used_codes.txt", "File containing used codes")
	flag.Parse()

	getCodeHandler, err := djv_codes.NewGetCodeHandler(*codesDir, *usedCodesPath)
	if err != nil {
		log.Fatal("Could not create handler: ", err)
	}

	http.Handle("/get", getCodeHandler)
	log.Fatal(http.ListenAndServe(*port, nil))
}
