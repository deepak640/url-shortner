package handlers
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"url-shortner/config"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func RemoveHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.DB.Database("urlshortener").Collection("urls")
	var body struct {
		Code   string `json:"code"`
		UserID string `json:"userid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if body.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if body.Code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := collection.FindOneAndDelete(ctx,bson.M{"short_code": body.Code,
		"user_id": body.UserID})



	if res.Err() != nil {
		http.Error(w, "No Link Found", http.StatusNotFound)
		return
	}

	log.Printf("Deleted short code: %s", body.Code)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Short URL removed"})
}
