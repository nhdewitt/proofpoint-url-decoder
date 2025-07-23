package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nhdewitt/proofpoint-url-decoder/internal/config"
)

// Decoder interface defines a method for decoding URLs.
// It is used by the APIDecodeHandler to decode URLs from the request.
// Implementations of this interface should provide the Decode method.
type Decoder interface {
	Decode(string) (string, error)
}

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
