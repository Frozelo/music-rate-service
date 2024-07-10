package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.ListenAndServe("localhost:8080", routes())
}

func routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/health", healthHandler)

	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
