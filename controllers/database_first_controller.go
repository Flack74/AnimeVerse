package controller

import (
	"encoding/json"
	"net/http"

	"animeverse/services"
	"github.com/go-chi/chi/v5"
)

// GetAnimeWithFallbackHandler handles database-first anime requests
func GetAnimeWithFallbackHandler(w http.ResponseWriter, r *http.Request) {
	animeName := chi.URLParam(r, "name")
	
	if animeName == "" {
		http.Error(w, "Anime name is required", http.StatusBadRequest)
		return
	}

	anime, err := services.GetAnimeWithFallback(animeName)
	if err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    anime,
	})
}