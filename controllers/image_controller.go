package controller

import (
	"encoding/json"
	"net/http"

	"animeverse/services"
)

func GetHighQualityImagesHandler(w http.ResponseWriter, r *http.Request) {
	animeName := r.URL.Query().Get("name")
	
	if animeName == "" {
		http.Error(w, "Anime name is required", http.StatusBadRequest)
		return
	}

	cover, banner, err := services.GetHighQualityImages(animeName)
	if err != nil {
		http.Error(w, "Failed to get high quality images", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"cover":   cover,
		"banner":  banner,
	})
}