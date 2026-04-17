package handlers

import (
	"context"
	"encoding/json"
	"os"

	"log"
	"strings"

	"math/rand"
	"net/http"
	"time"
	"url-shortner/config"
	"url-shortner/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)



func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.DB.Database("urlshortener").Collection("urls")
	var body struct {
		URL string `json:"url"`
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
	Server := os.Getenv("SERVER")
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"short_url": Server + code,
	}

	json.NewEncoder(w).Encode(response)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.Trim(r.URL.Path, "/")

	// ✅ Avoid empty or reserved paths
	if code == "" || code == "shorten" {
		http.Error(w, "Invalid short URL", http.StatusBadRequest)
		return
	}

	collection := config.DB.Database("urlshortener").Collection("urls")

	var result models.URL

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"short_code": code}).Decode(&result)

	// ✅ Proper error handling
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Short code not found: %s", code)
			http.NotFound(w, r)
			return
		}

		log.Printf("DB error for code %s: %v", code, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}


	// ✅ Ensure URL has scheme
	if !strings.HasPrefix(result.LongURL, "http") {
		result.LongURL = "http://" + result.LongURL
	}

	http.Redirect(w, r, result.LongURL, http.StatusFound)
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
