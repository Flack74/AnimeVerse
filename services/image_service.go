package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"animeverse/config"
	model "animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var imageCollection *mongo.Collection

// InitImageCollection initializes the image collection
func InitImageCollection() {
	if config.DB != nil {
		imageCollection = config.GetCollection(config.DB, "image_cache")
	}
}

// GetOrFetchImages gets images from cache or fetches from APIs
func GetOrFetchImages(malID, anilistID int) (*model.ImageCache, error) {
	// First check cache
	if cached := getCachedImages(malID, anilistID); cached != nil {
		// Check if cache is stale (30 days)
		if time.Since(cached.LastUpdated).Hours() < 24*30 {
			return cached, nil
		}
	}

	// Fetch from APIs
	imageURL := ""
	bannerURL := ""

	if malID > 0 {
		if img, err := fetchJikanImage(malID); err == nil {
			imageURL = img
		}
	}

	if anilistID > 0 {
		if banner, err := fetchAniListBanner(anilistID); err == nil {
			bannerURL = banner
		}
	}

	// Save to cache
	imageCache := &model.ImageCache{
		MALID:       malID,
		AniListID:   anilistID,
		ImageUrl:    imageURL,
		BannerUrl:   bannerURL,
		LastUpdated: time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := saveImageCache(imageCache); err != nil {
		return nil, err
	}

	return imageCache, nil
}

// getCachedImages retrieves images from cache
func getCachedImages(malID, anilistID int) *model.ImageCache {
	if imageCollection == nil {
		InitImageCollection()
	}
	if imageCollection == nil {
		return nil
	}
	
	var imageCache model.ImageCache
	filter := bson.M{}
	
	if malID > 0 {
		filter["mal_id"] = malID
	}
	if anilistID > 0 {
		filter["anilist_id"] = anilistID
	}

	if len(filter) == 0 {
		return nil
	}

	err := imageCollection.FindOne(context.Background(), filter).Decode(&imageCache)
	if err != nil {
		return nil
	}

	return &imageCache
}

// fetchJikanImage fetches image from Jikan API
func fetchJikanImage(malID int) (string, error) {
	url := fmt.Sprintf("https://api.jikan.moe/v4/anime/%d", malID)
	
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("jikan API returned status %d", resp.StatusCode)
	}

	var jikanResp model.JikanResponse
	if err := json.NewDecoder(resp.Body).Decode(&jikanResp); err != nil {
		return "", err
	}

	return jikanResp.Data.Images.JPG.LargeImageURL, nil
}

// fetchAniListBanner fetches banner from AniList GraphQL API
func fetchAniListBanner(anilistID int) (string, error) {
	query := `
	query ($id: Int) {
		Media (id: $id, type: ANIME) {
			bannerImage
			coverImage {
				large
			}
		}
	}`

	variables := map[string]interface{}{
		"id": anilistID,
	}

	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("aniList API returned status %d", resp.StatusCode)
	}

	var anilistResp model.AniListResponse
	if err := json.NewDecoder(resp.Body).Decode(&anilistResp); err != nil {
		return "", err
	}

	if anilistResp.Data.Media.BannerImage != "" {
		return anilistResp.Data.Media.BannerImage, nil
	}

	return anilistResp.Data.Media.CoverImage.Large, nil
}

// saveImageCache saves image cache to database
func saveImageCache(imageCache *model.ImageCache) error {
	if imageCollection == nil {
		InitImageCollection()
	}
	if imageCollection == nil {
		return fmt.Errorf("image collection not initialized")
	}
	
	filter := bson.M{}
	if imageCache.MALID > 0 {
		filter["mal_id"] = imageCache.MALID
	}
	if imageCache.AniListID > 0 {
		filter["anilist_id"] = imageCache.AniListID
	}

	update := bson.M{
		"$set": bson.M{
			"image_url":    imageCache.ImageUrl,
			"banner_url":   imageCache.BannerUrl,
			"last_updated": imageCache.LastUpdated,
		},
		"$setOnInsert": bson.M{
			"created_at": imageCache.CreatedAt,
		},
	}

	upsert := true
	_, err := imageCollection.UpdateOne(
		context.Background(),
		filter,
		update,
		&options.UpdateOptions{Upsert: &upsert},
	)

	return err
}

// SaveImageData saves image data from frontend request
func SaveImageData(req model.ImageRequest) error {
	imageCache := &model.ImageCache{
		MALID:       req.MALID,
		AniListID:   req.AniListID,
		ImageUrl:    req.ImageUrl,
		BannerUrl:   req.BannerUrl,
		LastUpdated: time.Now(),
		CreatedAt:   time.Now(),
	}

	return saveImageCache(imageCache)
}

// GetImagesByIDs gets cached images by MAL/AniList IDs
func GetImagesByIDs(malID, anilistID int) *model.ImageCache {
	if imageCollection == nil {
		InitImageCollection()
	}
	return getCachedImages(malID, anilistID)
}