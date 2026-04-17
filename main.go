package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"url-shortner/config"
)


type shortner struct {
	URL string `json:"url"`
}


func main() {
	config.Connect()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	mux.HandleFunc("POST /shorten", shortenURL)

	fmt.Println("Server started on :8080")
	err:= http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}


// handlers
func shortenURL(w http.ResponseWriter, r *http.Request) {

	var body shortner

	err:= json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received URL to shorten: %s\n", body.URL)
	json.NewEncoder(w).Encode(map[string]string{
		"shortened_url": body.URL,
	})
}


func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
