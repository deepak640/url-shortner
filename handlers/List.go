package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"url-shortner/config"
	"url-shortner/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.DB.Database("urlshortener").Collection("urls")
	var body struct {
		UserID string `json:"userid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	res,err := collection.Find(context.TODO(), bson.M{"user_id": body.UserID })

	if err != nil {
		http.Error(w,"Links Not Found for this User",http.StatusNotFound)
	}
	var urls []models.URL
	if err := res.All(context.TODO(), &urls); err != nil {
		http.Error(w, "Error parsing URLs", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(urls)
}
