package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
	"url-shortner/config"
	"url-shortner/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

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
