package controller

import (
	"encoding/json"
	"net/http"

	"animeverse/services"
)

func GetSpotlightHandler(w http.ResponseWriter, r *http.Request) {
	spotlight, err := services.GetSpotlightAnime()
	if err != nil {
		http.Error(w, "Failed to get spotlight anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    spotlight,
	})
}

func GetTopRatedMixedHandler(w http.ResponseWriter, r *http.Request) {
	topRated, err := services.GetTopRatedMixed()
	if err != nil {
		http.Error(w, "Failed to get top rated anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    topRated,
	})
}