package controller

import (
	"encoding/json"
	"net/http"

	"animeverse/services"
)

func UpgradeImagesHandler(w http.ResponseWriter, r *http.Request) {
	animeName := r.URL.Query().Get("name")
	animeID := r.URL.Query().Get("id")
	currentBanner := r.URL.Query().Get("currentBanner")
	currentCover := r.URL.Query().Get("currentCover")
	
	if animeName == "" {
		http.Error(w, "Anime name is required", http.StatusBadRequest)
		return
	}

	newBanner, newCover, err := services.UpgradeImageQualityAndSave(animeID, animeName, currentBanner, currentCover)
	if err != nil {
		http.Error(w, "Failed to upgrade images", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"banner":  newBanner,
		"cover":   newCover,
		"upgraded": newBanner != currentBanner || newCover != currentCover,
	})
}