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

	// ✅ Check if link is active
	if !result.IsActive {
		http.Error(w, "Link is no longer active", http.StatusGone)
		return
	}

	// ✅ Check if link is expired
	if result.ExpiresAt != nil && result.ExpiresAt.Before(time.Now().UTC()) {
		collection.UpdateOne(ctx, bson.M{"_id": result.ID}, bson.M{"$set": bson.M{"is_active": false}})
		log.Printf("Short code expired: %s", code)
		http.Error(w, "Link has expired", http.StatusGone)
		return
	}

	// ✅ Check if max clicks reached
	if result.MaxClicks > 0 && result.CurrentClicks >= result.MaxClicks {
		collection.UpdateOne(ctx, bson.M{"_id": result.ID}, bson.M{"$set": bson.M{"is_active": false}})
		log.Printf("Max clicks reached for code %s", code)
		http.Error(w, "Link has reached maximum usage", http.StatusGone)
		return
	}

	// ✅ Increment click count
	_, err = collection.UpdateOne(ctx, bson.M{"_id": result.ID}, bson.M{"$inc": bson.M{"current_clicks": 1}})
	if err != nil {
		log.Printf("Error updating click count for %s: %v", code, err)
	}

	// ✅ Ensure URL has scheme
	if !strings.HasPrefix(result.LongURL, "http") {
		result.LongURL = "http://" + result.LongURL
	}

	http.Redirect(w, r, result.LongURL, http.StatusFound)
}
