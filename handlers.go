package main

import (
	"encoding/json"
	"html"
	"html/template"
	"net/http"
	"path"
	"strings"
)

//	FileServer serves the static HTML, CSS, and JS files. Redirects to index if user tries to access a directory or file without an extension.
func FileServer(fs http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/static/")

		switch p {
		case "", "index.html", "js", "css":
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if strings.HasSuffix(p, "/") || path.Ext(p) == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		fs.ServeHTTP(w, r)
	})
}

//	FormHandler serves the main form page. Rejects any request other than GET.
// 	It uses a template located in templates/form.html
func FormHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "form", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

//	APIDecodeHandler handles POST requests to decode URLs via JSON API.
//	It expects a JSON body with a list of URLs and returns a JSON response with decoded results or errors.
//	It uses a Decoder interface to decode the URLs.
func APIDecodeHandler(dec Decoder) http.HandlerFunc {
	type request struct {
		URLs []string `json:"urls"`
	}
	type response struct {
		Results []string `json:"results"`
		Errors  []string `json:"errors,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		resp := response{
			Results: make([]string, len(req.URLs)),
			Errors:  make([]string, len(req.URLs)),
		}
		for i, u := range req.URLs {
			u = html.UnescapeString(u)
			out, err := dec.Decode(u)
			resp.Results[i] = out
			if err != nil {
				resp.Errors[i] = err.Error()
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

//	DecodeFormHandler handles form submissions for decoding URLs.
//	It expects a POST request with a form field named "input" containing URLs.
//	It decodes each URL and returns the results in a template.
func DecodeFormHandler(dec Decoder, tmpl *template.Template) http.HandlerFunc {
	type result struct {
		Input  string
		Output string
		Err    error
	}

	return func(w http.ResponseWriter, r *http.Request) {
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

		var results []result
		for _, line := range lines {
			u := strings.TrimSpace(line)
			if u == "" {
				continue
			}
			out, err := dec.Decode(u)
			results = append(results, result{
				Input:  u,
				Output: out,
				Err:    err,
			})
		}

		if err := tmpl.ExecuteTemplate(w, "result", results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func MobileFormHandler(dec Decoder, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.ExecuteTemplate(w, "mobile_form", nil)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		raw := strings.TrimSpace(r.FormValue("input"))
		if raw == "" {
			http.Error(w, "empty input", http.StatusBadRequest)
			return
		}

		decoded, err := dec.Decode(html.UnescapeString(raw))
		data := struct {
			Input  string
			Output string
			Err    error
		}{
			Input:  raw,
			Output: decoded,
			Err:    err,
		}

		tmpl.ExecuteTemplate(w, "mobile_result", data)
	}
}

func MobileRedirectHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if strings.Contains(ua, "Mobile") &&
			!strings.HasPrefix(r.URL.Path, "/m") &&
			r.Method == http.MethodGet {
			http.Redirect(w, r, "/m", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
