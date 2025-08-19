package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to product-service")
	})
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"id": 101, "name": "Laptop"}, {"id": 102, "name": "Phone"}]`)
	})
	fmt.Println("Product service running on :8080")
	http.ListenAndServe(":8080", nil)
}
