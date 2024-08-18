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

	databasePath := "database.json"
	db, err := NewDB(databasePath)
	if err != nil {
		log.Fatal("Could not connect to a database")
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))

	fsHandler := http.StripPrefix("/app", fs)
	mux.Handle("/app/*", apiState.middlewareMetricsInc(fsHandler))

	mux.HandleFunc("GET /admin/metrics", apiState.handlerMetricsCount)
	mux.HandleFunc("GET /api/reset", apiState.handlerMetricsReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	go mux.HandleFunc("POST /api/chirps", apiState.handlerCreateChirp(db))
	go mux.HandleFunc("GET /api/chirps", apiState.handlerGetChirps(db))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Running server at: http://%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
