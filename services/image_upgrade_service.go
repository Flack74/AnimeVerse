package services

import (
	"context"
	"time"

	"animeverse/cache"
	"animeverse/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpgradeImageQualityAndSave(animeID, animeName string, currentBanner, currentCover string) (string, string, error) {
	// Check if images need upgrading
	if !needsImageUpgrade(currentBanner, currentCover) {
		return currentBanner, currentCover, nil
	}

	// Get high-quality images from AniList
	newCover, newBanner, err := GetHighQualityImages(animeName)
	if err != nil {
		return currentBanner, currentCover, err
	}

	// Update database if we have an ID
	if animeID != "" {
		updateImageInDatabase(animeID, newCover, newBanner)
	}

	return newBanner, newCover, nil
}

func needsImageUpgrade(banner, cover string) bool {
	return isLowQualityImage(banner) || isLowQualityImage(cover)
}

func isLowQualityImage(imageUrl string) bool {
	if imageUrl == "" {
		return true
	}
	
	lowQualityIndicators := []string{
		"placeholder", "small", "thumb", "150x", "300x", 
		"medium", "default", "no-image", "missing",
	}
	
	for _, indicator := range lowQualityIndicators {
		if contains(imageUrl, indicator) {
			return true
		}
	}
	
	return false
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    (len(str) > len(substr) && 
		     (str[:len(substr)] == substr || 
		      str[len(str)-len(substr):] == substr ||
		      indexOf(str, substr) >= 0)))
}

func indexOf(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func updateImageInDatabase(animeID, newCover, newBanner string) error {
	objectID, err := primitive.ObjectIDFromHex(animeID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"imageUrl":   newCover,
			"bannerUrl":  newBanner,
			"updated_at": time.Now(),
		},
	}

	_, err = config.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)

	// Clear related caches
	if err == nil {
		clearImageCaches(animeID)
	}

	return err
}

func clearImageCaches(animeID string) {
	// Clear anime-specific caches
	cache.Delete("anime:" + animeID)
	cache.Delete("enhanced_anime:" + animeID)
	
	// Clear general caches that might contain this anime
	cache.Delete("spotlight:anime")
	cache.Delete("trending:anime")
	cache.Delete("popular:anime")
}

func UpgradeSpotlightImages() error {
	// Get spotlight anime from database
	spotlight, err := getSpotlightFromDatabase()
	if err != nil {
		return err
	}

	for _, anime := range spotlight {
		if needsImageUpgrade(anime.BannerUrl, anime.ImageUrl) {
			newCover, newBanner, err := GetHighQualityImages(anime.Name)
			if err == nil {
				updateImageInDatabase(anime.ID, newCover, newBanner)
			}
		}
	}

	return nil
}