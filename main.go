package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func main() {
	const port = "8080"
	addr := "localhost:" + port
	apiState := apiConfig{
		fileServerHits: 0,
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))
	fsHandler := http.StripPrefix("/app", fs)
	mux.Handle("/app/*", apiState.middlewareMetricsInc(fsHandler))

	mux.Handle("/metrics", apiState.middlewareMetricsCount())
	mux.Handle("/reset", apiState.middlewareMetricsReset())
	mux.HandleFunc("/healthz", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Running server at: http://%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) middlewareMetricsReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits = 0
	})
}

func (cfg *apiConfig) middlewareMetricsCount() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits)))
	})
}
