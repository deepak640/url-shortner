package handlers

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"url-shortner/config"
	"url-shortner/models"
	"github.com/joho/godotenv"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	Server := os.Getenv("SERVER")
	collection := config.DB.Database("urlshortener").Collection("urls")
	var body struct {
		URL    string `json:"url"`
		UserID string `json:"userid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if body.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	code := generateCode()

	doc := models.URL{
		ShortCode: code,
		UserID:    body.UserID,
		LongURL:   body.URL,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), doc)

	if err != nil {
		log.Printf("Error saving URL to DB: %v", err)
		http.Error(w, "Error saving URL", 500)
		return
	}

	log.Printf("Shortened %s to %s", body.URL, code)


	response := map[string]string{
		"short_url": Server + code,
	}

	json.NewEncoder(w).Encode(response)
}


// functions
func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
