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

	manager, err := djv_codes.NewCodeManager(*codesDir, *usedCodesPath)
	if err != nil {
		log.Fatalf("Couldn't create code manager %v", err)
	}

	http.Handle("/get", &djv_codes.GetCodeHandler{manager})
	http.Handle("/list", &djv_codes.ListCodeHandler{manager})
	log.Fatal(http.ListenAndServe(*port, nil))
}
