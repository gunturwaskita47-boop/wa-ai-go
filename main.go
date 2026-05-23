package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "WhatsApp AI Bot Online")
	})

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Webhook Active")
	})

	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, nil)
}
