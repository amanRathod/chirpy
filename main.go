package main

import (
	"html/template"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Responsible for routing Http requests
	mux := http.NewServeMux()

	// File server handler with middleware
	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", apiCfg.middlewareMetricsIncrement(http.StripPrefix("/app",fileServer)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

	// http struct
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}


// Middleware to increment the fileserver hits counter
func (apiCfg *apiConfig) middlewareMetricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)        // Call the next handler
	})
}

// Handler to display the current metrics
func (apiCfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Increment the file server hits counter
	hits := apiCfg.fileserverHits.Load()

	// Parse and execute the template
	tmpl, err := template.ParseFiles("./metrics.html")
	if err != nil {
		http.Error(w, "Could not parse metrics.html", http.StatusInternalServerError)
		return
	}

	data := struct {
		Hits int32
	}{
		Hits: hits,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}

}