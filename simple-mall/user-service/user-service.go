package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to user-service")
	})
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"id": 1, "name": "Alice"}`)
	})
	fmt.Println("User service running on :8080")
	http.ListenAndServe(":8080", nil)
}
