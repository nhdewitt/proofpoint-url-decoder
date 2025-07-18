package main

import (
	"html"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.New("page").Parse(`
{{define "form"}}
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Proofpoint URL Decoder</title></head>
<body>
	<h1>Proofpoint URL Decoder</h1>
	<form method="POST" action="/decode">
		<input type="text" name="input" placeholder="Enter encoded URL" size="50">
		<button type="submit">Decode</button>
	</form>
</body>
</html>
{{end}}

{{define "result"}}
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Decoded URL</title></head>
<body>
	<h1>Decoded URL</h1>
	<p><code>{{.}}</code></p>
	<a href="/">Decode another</a>
</body>
</html>
{{end}}
`))

func runServer(d *urlDefenseDecoder) {
	mux := http.NewServeMux()

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
		raw := r.FormValue("input")

		u := html.UnescapeString(raw)
		decoded, err := d.Decode(u)
		if err != nil {
			http.Error(w, "Invalid URL encoding", http.StatusBadRequest)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "result", decoded); err != nil {
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
