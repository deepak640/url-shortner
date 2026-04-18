package main

import (
	"fmt"
	"net/http"
	"os"
	"url-shortner/config"
	"url-shortner/handlers"
	middleware "url-shortner/middlware"
)




func main() {
	config.Connect()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", handlers.ShortenHandler)
	mux.HandleFunc("/",handlers.RedirectHandler)
	mux.HandleFunc("POST /list",handlers.ListHandler)
	mux.HandleFunc("POST /remove",handlers.RemoveHandler)
	logger := middleware.Logger(mux)
	port := os.Getenv("PORT")
    if port == "" {
        port = "3000" // fallback for local dev only
    }
	fmt.Println("Server started on :" + port)
	err:= http.ListenAndServe(":" + port, logger)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
