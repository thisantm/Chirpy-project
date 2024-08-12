package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	addr := "localhost:" + port

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Running server at: http://%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
