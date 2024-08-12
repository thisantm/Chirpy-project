package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := ":8080"
	addr := "localhost" + port

	http.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Printf("Running server at: http://%s\n", addr)
	http.ListenAndServe(port, nil)
}
