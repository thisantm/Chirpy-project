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
	fsAdminHandler := http.StripPrefix("/admin", fs)
	fsHandler := http.StripPrefix("/app", fs)
	mux.Handle("/app/*", apiState.middlewareMetricsInc(fsHandler))

	mux.Handle("GET /admin/metrics", apiState.middlewareMetricsCount(fsAdminHandler))
	mux.Handle("GET /api/reset", apiState.middlewareMetricsReset())
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Running server at: http://%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
