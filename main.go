package main

import (
	"fmt"
	"net/http"
)

func main() {
	addr := "localhost:8080"
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	fmt.Printf("Running server at: http://%s\n", addr)
	srv.ListenAndServe()
}
