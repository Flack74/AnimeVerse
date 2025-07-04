package controller

import (
	"encoding/json"
	"net/http"

	"animeverse/services"
)

func GetAnimeThemesHandler(w http.ResponseWriter, r *http.Request) {
	animeName := r.URL.Query().Get("name")
	
	if animeName == "" {
		http.Error(w, "Anime name is required", http.StatusBadRequest)
		return
	}

	themes, err := services.GetAnimeThemes(animeName)
	if err != nil {
		http.Error(w, "Failed to get themes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    themes,
	})
}