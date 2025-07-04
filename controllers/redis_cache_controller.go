package controller

import (
	"encoding/json"
	"net/http"

	"animeverse/services"
)

func GetTrendingFastHandler(w http.ResponseWriter, r *http.Request) {
	animes, err := services.GetTrendingWithRedisCache()
	if err != nil {
		http.Error(w, "Failed to get trending anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    animes,
	})
}