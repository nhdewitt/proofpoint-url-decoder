package main

import (
	"encoding/json"
	"html"
	"html/template"
	"net/http"
	"strings"
)

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
