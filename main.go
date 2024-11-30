package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/thisantm/Chirpy-project/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to postgres database")
	}

	dbQueries := database.New(db)

	const port = "8080"
	addr := "localhost:" + port
	apiState := apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries:      dbQueries,
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))

	fsHandler := http.StripPrefix("/app", fs)
	mux.Handle("/app/*", apiState.middlewareMetricsInc(fsHandler))

	mux.HandleFunc("GET /admin/metrics", apiState.handlerMetricsCount)
	mux.HandleFunc("GET /api/reset", apiState.handlerMetricsReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", apiState.handlerValidateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Running server at: http://%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
