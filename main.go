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

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/app/*", http.StripPrefix("/app", fs))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/healthz" {
			http.NotFound(w, req)
			return
		}
		body := []byte("OK")
		req.Header.Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write(body)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Running server at: http://%s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
