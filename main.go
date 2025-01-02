package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/thisantm/Chirpy-project/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	platform       string
	db             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to postgres database")
	}

	dbQueries := database.New(db)

	const port = "8080"
	addr := "localhost:" + port
	apiState := apiConfig{
		fileServerHits: atomic.Int32{},
		platform:       platform,
		db:             dbQueries,
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))

	fsHandler := http.StripPrefix("/app", fs)
	mux.Handle("/app/*", apiState.middlewareMetricsInc(fsHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", apiState.handlerUsersCreate)
	mux.HandleFunc("POST /api/chirps", apiState.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiState.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiState.handlerGetChirp)

	mux.HandleFunc("GET /admin/metrics", apiState.handlerMetricsCount)
	mux.HandleFunc("POST /admin/reset", apiState.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Running server at: http://%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
