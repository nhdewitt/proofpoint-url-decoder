package main

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func runServer(d *urlDefenseDecoder) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "form", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		raw := html.UnescapeString(r.FormValue("input"))

		lines := strings.Split(raw, "\n")

		type result struct {
			Input  string
			Output string
			Err    error
		}
		var results []result
		for _, line := range lines {
			u := strings.TrimSpace(line)
			if u == "" {
				continue
			}
			decoded, err := d.Decode(u)
			results = append(results, result{
				Input:  u,
				Output: decoded,
				Err:    err,
			})
		}

		if err := tmpl.ExecuteTemplate(w, "result", results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Listening on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
