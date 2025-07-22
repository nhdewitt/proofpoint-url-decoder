package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nhdewitt/proofpoint-url-decoder/internal/config"
)

type Decoder interface {
	Decode(string) (string, error)
}

func runServer(d Decoder, c config.Config) {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", FormHandler(tmpl))
	mux.HandleFunc("/api/decode", APIDecodeHandler(d))
	mux.HandleFunc("/decode", DecodeFormHandler(d, tmpl))

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + c.Port,
	}

	log.Printf("Listening on port: %s\n", c.Port)
	log.Fatal(server.ListenAndServe())
}
