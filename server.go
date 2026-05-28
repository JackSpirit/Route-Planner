package main

import (
	"fmt"
	"net/http"
)

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/route", handleRoutePlaceholder)

	fmt.Println("Web server starting on http://localhost:8080..")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

func handleRoutePlaceholder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Router working hopefully :)"}`))
}