package main

import (
	"fmt"
	"net/http"
	"url-shortner/config"
	"url-shortner/handlers"
	middleware "url-shortner/middlware"
)




func main() {
	config.Connect()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", handlers.ShortenHandler)
	mux.HandleFunc("/",handlers.RedirectHandler)
	logger := middleware.Logger(mux)
	fmt.Println("Server started on :8080")
	err:= http.ListenAndServe(":8080", logger)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
