package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Bazel + Go!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
