package handlers

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"url-shortner/config"
	"url-shortner/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	Server := os.Getenv("SERVER")
	collection := config.DB.Database("urlshortener").Collection("urls")
	var body struct {
		URL       string `json:"URL"`
		UserID    string `json:"userid"`
		ExpiresIn string `json:"ExpiresIn"`
		CustomCode string `json:"CustomCode"`
		MaxClicks string `json:"MaxClicks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if body.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if(body.ExpiresIn == ""){
		body.ExpiresIn = "24h"
	}

	if(body.MaxClicks == ""){
		body.MaxClicks = "100"
	}
	// Parse Expiry
	expiryDate := parseExpiryDuration(body.ExpiresIn)

	// Parse MaxClicks
	maxClicks, _ := strconv.Atoi(body.MaxClicks)
	response := collection.FindOne(context.TODO(), bson.M{"short_code": body.CustomCode})
	if response.Err() == nil {
		http.Error(w, "Custom code already exists", http.StatusBadRequest)
		return
	}
	code := body.CustomCode
	if code == "" {
		code = generateCode()
	}

	doc := models.URL{
		ShortCode:     code,
		UserID:        body.UserID,
		LongURL:       body.URL,
		CreatedAt:     time.Now(),
		ExpiresAt:     expiryDate,
		MaxClicks:     maxClicks,
		CurrentClicks: 0,
		IsActive:      true,
	}

	_, err := collection.InsertOne(context.TODO(), doc)

	if err != nil {
		log.Printf("Error saving URL to DB: %v", err)
		http.Error(w, "Error saving URL", 500)
		return
	}

	log.Printf("Shortened %s to %s", body.URL, code)


	response1 := map[string]string{
		"short_url": Server + code,
	}

	json.NewEncoder(w).Encode(response1)
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

// Parse duration strings like "24h", "7d", "30m"
func parseExpiryDuration(expiresIn string) *time.Time {
	if expiresIn == "" {
		return nil
	}

	// Regular expression to match number + unit
	re := regexp.MustCompile(`^(\d+)([hHdDmM])?$`)
	matches := re.FindStringSubmatch(expiresIn)

	if len(matches) < 2 {
		return nil
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}

	unit := "h" // default to hours
	if len(matches) == 3 && matches[2] != "" {
		unit = strings.ToLower(matches[2])
	}

	var duration time.Duration

	switch unit {
	case "h":
		duration = time.Duration(value) * time.Hour
	case "d":
		duration = time.Duration(value) * 24 * time.Hour
	case "m":
		duration = time.Duration(value) * 30 * 24 * time.Hour // months
	default:
		return nil
	}

	expiryTime := time.Now().Add(duration)
	return &expiryTime
}
