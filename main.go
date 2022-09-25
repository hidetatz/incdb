package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", postQuery)

	s := http.Server{
		Addr:    ":2134",
		Handler: mux,
	}

	s.ListenAndServe()
}
