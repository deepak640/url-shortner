package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"url-shortner/config"
	"url-shortner/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.DB.Database("urlshortener").Collection("urls")
	var body struct {
		ShortCode string `json:"shortCode"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx,cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result models.URL
	err2 := collection.FindOne(ctx,bson.M{"short_code":body.ShortCode}).Decode(&result)
	if err2 != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(result)
}
