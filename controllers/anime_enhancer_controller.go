package controller

import (
	"encoding/json"
	"net/http"

	"animeverse/services"
	"github.com/go-chi/chi/v5"
)

type EnhanceRequest struct {
	Name string `json:"name"`
}

// EnhanceAnime enhances existing anime with API data
func EnhanceAnime(w http.ResponseWriter, r *http.Request) {
	animeID := chi.URLParam(r, "id")
	
	var req EnhanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := services.UpdateAnimeWithAPIData(animeID, req.Name)
	if err != nil {
		http.Error(w, "Failed to enhance anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Anime enhanced successfully",
	})
}

// CreateAnimeFromAPI creates new anime from API data
func CreateAnimeFromAPI(w http.ResponseWriter, r *http.Request) {
	var req EnhanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	anime, err := services.EnhanceAnimeFromAPI(req.Name)
	if err != nil {
		http.Error(w, "Failed to create anime from API", http.StatusInternalServerError)
		return
	}

	if anime == nil {
		http.Error(w, "Anime not found in API", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    anime,
		"message": "Anime created successfully",
	})
}

// GetEnhancedAnime gets anime with full details from external API
func GetEnhancedAnime(w http.ResponseWriter, r *http.Request) {
	animeID := chi.URLParam(r, "id")
	animeName := r.URL.Query().Get("name")
	
	if animeName == "" {
		http.Error(w, "Anime name is required", http.StatusBadRequest)
		return
	}

	anime, err := services.EnhanceAnimeWithFullData(animeID, animeName)
	if err != nil {
		http.Error(w, "Failed to enhance anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    anime,
	})
}