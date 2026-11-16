package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	listenAddress  = "127.0.0.1:8080"
	maxQueryLength = 100
)

type PageData struct {
	Query  string
	Groups []LinkGroup
	Error  string
}

func main() {
	links, warnings, err := loadLinks("links.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, warning := range warnings {
		log.Printf("configuration warning: %v (entry skipped)", warning)
	}
	log.Printf("loaded %d search links", len(links))

	pageTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("load HTML template: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := strings.TrimSpace(r.URL.Query().Get("q"))
		data := PageData{Query: query}
		status := http.StatusOK

		if utf8.RuneCountInString(query) > maxQueryLength {
			data.Error = "Search terms must be 100 characters or fewer."
			status = http.StatusBadRequest
		} else if query != "" {
			data.Groups = groupLinks(links, query)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		if err := pageTemplate.Execute(w, data); err != nil {
			log.Printf("render page: %v", err)
		}
	})

	server := &http.Server{
		Addr:              listenAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("dictionary search is available at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server: %v", err)
	}
}
