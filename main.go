package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/thisantm/Chirpy-project/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits int
	dbQueries      *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	const port = "8080"
	addr := "localhost:" + port
	apiState := apiConfig{
		fileServerHits: 0,
		dbQueries:      dbQueries,
	}

	// databasePath := "database.json"
	// db, err := NewDB(databasePath)
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
	// go mux.HandleFunc("POST /api/chirps", apiState.handlerCreateChirp(db))
	// go mux.HandleFunc("GET /api/chirps/", apiState.handlerGetChirps(db))
	// go mux.HandleFunc("GET /api/chirps/{id}", apiState.handlerGetChirpById(db))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Running server at: http://%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
