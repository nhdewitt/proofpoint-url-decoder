package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nhdewitt/proofpoint-url-decoder/internal/config"
)

// Decoder defines an interface for decoding encoded URLs.
// Implementations should return the decoded URL or an error.
type Decoder interface {
	Decode(string) (string, error)
}

// runServer starts an HTTP server that provides a web interface for decoding
// Proofpoint URL Defense links.
//
// It serves templates from the `templates/` directory and uses the provided
// Decoder implementation to process input URLs.
func runServer(d Decoder, c config.Config) {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	mux := http.NewServeMux()

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	mux.Handle("/static/", FileServer(fileServer))

	mux.Handle("/", MobileRedirectHandler(http.HandlerFunc(FormHandler(tmpl))))
	mux.HandleFunc("/api/decode", APIDecodeHandler(d))
	mux.Handle("/decode", MobileRedirectHandler(http.HandlerFunc(DecodeFormHandler(d, tmpl))))
	mux.HandleFunc("/m", MobileFormHandler(d, tmpl))

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + c.Port,
	}

	log.Printf("Listening on port: %s\n", c.Port)
	log.Fatal(server.ListenAndServe())
}
