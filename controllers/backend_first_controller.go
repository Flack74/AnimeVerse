package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"animeverse/cache"
	"animeverse/config"
	"animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BackendFirstTrendingHandler - Backend first with high-quality images
func BackendFirstTrendingHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	cacheKey := fmt.Sprintf("backend_trending_%d", page)
	
	// Check Redis cache first
	var cachedAnimes []models.Anime
	if cache.Exists(cacheKey) {
		if err := cache.Get(cacheKey, &cachedAnimes); err == nil {
			log.Printf("✅ Cache hit for trending page %d", page)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    cachedAnimes,
				"source":  "cache",
			})
			return
		}
	}

	// Check MongoDB for existing data
	collection := config.GetCollection(config.DB, "anime")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (page - 1) * 24
	opts := options.Find().SetSkip(int64(skip)).SetLimit(24).SetSort(bson.D{{"score", -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	
	var dbAnimes []models.Anime
	if err == nil {
		cursor.All(ctx, &dbAnimes)
		cursor.Close(ctx)
	}

	// If we have good quality data in DB, use it
	if len(dbAnimes) >= 20 {
		log.Printf("✅ Using MongoDB data for trending page %d", page)
		
		// Upgrade images in background
		go upgradeAnimeImages(dbAnimes)
		
		// Cache for 10 minutes
		cache.Set(cacheKey, dbAnimes, 10*time.Minute)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    dbAnimes,
			"source":  "database",
		})
		return
	}

	// Fallback to AniList with high-quality data
	log.Printf("🌐 Fetching from AniList for trending page %d", page)
	animes, err := fetchTrendingFromAniList(page)
	if err != nil {
		http.Error(w, "Failed to load trending anime", http.StatusInternalServerError)
		return
	}

	// Save to MongoDB in background
	go saveAnimesToDB(animes)
	
	// Cache for 5 minutes
	cache.Set(cacheKey, animes, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    animes,
		"source":  "anilist",
	})
}

// BackendFirstBrowseHandler - Backend first browse with filters
func BackendFirstBrowseHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	genre := r.URL.Query().Get("genre")
	year := r.URL.Query().Get("year")
	search := r.URL.Query().Get("search")

	cacheKey := fmt.Sprintf("backend_browse_%d_%s_%s_%s", page, genre, year, search)
	
	// Check cache first
	var cachedAnimes []models.Anime
	if cache.Exists(cacheKey) {
		if err := cache.Get(cacheKey, &cachedAnimes); err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    cachedAnimes,
				"source":  "cache",
			})
			return
		}
	}

	// Build MongoDB query
	filter := bson.M{}
	if genre != "" {
		filter["genre"] = bson.M{"$in": []string{genre}}
	}
	if year != "" {
		if yearInt, err := strconv.Atoi(year); err == nil {
			filter["year"] = yearInt
		}
	}
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	// Try MongoDB first
	collection := config.GetCollection(config.DB, "anime")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (page - 1) * 25
	opts := options.Find().SetSkip(int64(skip)).SetLimit(25).SetSort(bson.D{{"score", -1}})
	cursor, err := collection.Find(ctx, filter, opts)
	
	var dbAnimes []models.Anime
	if err == nil {
		cursor.All(ctx, &dbAnimes)
		cursor.Close(ctx)
	}

	// If we have enough data, use it
	if len(dbAnimes) >= 15 {
		go upgradeAnimeImages(dbAnimes)
		cache.Set(cacheKey, dbAnimes, 15*time.Minute)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    dbAnimes,
			"source":  "database",
		})
		return
	}

	// Fallback to AniList
	animes, err := fetchBrowseFromAniList(page, genre, year, search)
	if err != nil {
		http.Error(w, "Failed to load anime", http.StatusInternalServerError)
		return
	}

	go saveAnimesToDB(animes)
	cache.Set(cacheKey, animes, 10*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    animes,
		"source":  "anilist",
	})
}

// fetchTrendingFromAniList fetches trending anime with high-quality images
func fetchTrendingFromAniList(page int) ([]models.Anime, error) {
	query := `{
		Page(page: %d, perPage: 24) {
			media(type: ANIME, sort: TRENDING_DESC) {
				id
				title { romaji english native }
				description(asHtml: false)
				genres
				episodes
				status
				format
				startDate { year }
				averageScore
				coverImage { extraLarge large medium }
				bannerImage
				studios { nodes { name } }
			}
		}
	}`

	requestBody := map[string]interface{}{
		"query": fmt.Sprintf(query, page),
	}

	jsonData, _ := json.Marshal(requestBody)
	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return processAniListResponse(result)
}

// fetchBrowseFromAniList fetches browse anime with filters
func fetchBrowseFromAniList(page int, genre, year, search string) ([]models.Anime, error) {
	query := `{
		Page(page: %d, perPage: 25) {
			media(type: ANIME, sort: POPULARITY_DESC%s) {
				id
				title { romaji english native }
				description(asHtml: false)
				genres
				episodes
				status
				format
				startDate { year }
				averageScore
				coverImage { extraLarge large medium }
				bannerImage
				studios { nodes { name } }
			}
		}
	}`

	filters := ""
	if search != "" {
		filters += fmt.Sprintf(`, search: "%s"`, search)
	}
	if genre != "" {
		filters += fmt.Sprintf(`, genre_in: ["%s"]`, genre)
	}
	if year != "" {
		filters += fmt.Sprintf(`, seasonYear: %s`, year)
	}

	requestBody := map[string]interface{}{
		"query": fmt.Sprintf(query, page, filters),
	}

	jsonData, _ := json.Marshal(requestBody)
	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return processAniListResponse(result)
}

// processAniListResponse converts AniList response to our models
func processAniListResponse(result map[string]interface{}) ([]models.Anime, error) {
	var animes []models.Anime
	
	if data, ok := result["data"].(map[string]interface{}); ok {
		if page, ok := data["Page"].(map[string]interface{}); ok {
			if media, ok := page["media"].([]interface{}); ok {
				for _, item := range media {
					if anime, ok := item.(map[string]interface{}); ok {
						processedAnime := convertToAnimeModel(anime)
						animes = append(animes, processedAnime)
					}
				}
			}
		}
	}
	
	return animes, nil
}

// convertToAnimeModel converts AniList data to our Anime model
func convertToAnimeModel(data map[string]interface{}) models.Anime {
	anime := models.Anime{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Title
	if title, ok := data["title"].(map[string]interface{}); ok {
		if english, ok := title["english"].(string); ok && english != "" {
			anime.Name = english
		} else if romaji, ok := title["romaji"].(string); ok {
			anime.Name = romaji
		}
	}

	// Basic info
	if desc, ok := data["description"].(string); ok {
		anime.Synopsis = desc
	}
	if genres, ok := data["genres"].([]interface{}); ok {
		for _, g := range genres {
			if genre, ok := g.(string); ok {
				anime.Genre = append(anime.Genre, genre)
			}
		}
	}
	if year, ok := data["startDate"].(map[string]interface{}); ok {
		if y, ok := year["year"].(float64); ok {
			anime.Year = int(y)
		}
	}
	if format, ok := data["format"].(string); ok {
		anime.Type = models.AnimeType(format)
	}
	if score, ok := data["averageScore"].(float64); ok {
		anime.Score = score / 10.0
	}

	// High-quality images - prioritize extraLarge
	if coverImage, ok := data["coverImage"].(map[string]interface{}); ok {
		if extraLarge, ok := coverImage["extraLarge"].(string); ok && extraLarge != "" {
			anime.ImageUrl = extraLarge
		} else if large, ok := coverImage["large"].(string); ok && large != "" {
			anime.ImageUrl = large
		}
	}
	if banner, ok := data["bannerImage"].(string); ok {
		anime.BannerUrl = banner
	}

	// Information
	info := models.AnimeInformation{}
	if episodes, ok := data["episodes"].(float64); ok {
		info.Episodes = int(episodes)
	}
	if status, ok := data["status"].(string); ok {
		info.Status = status
	}
	if studios, ok := data["studios"].(map[string]interface{}); ok {
		if nodes, ok := studios["nodes"].([]interface{}); ok {
			for _, node := range nodes {
				if studio, ok := node.(map[string]interface{}); ok {
					if name, ok := studio["name"].(string); ok {
						info.Studios = append(info.Studios, name)
					}
				}
			}
		}
	}
	anime.Information = info

	return anime
}

// saveAnimesToDB saves anime data to MongoDB
func saveAnimesToDB(animes []models.Anime) {
	collection := config.GetCollection(config.DB, "anime")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, anime := range animes {
		filter := bson.M{"name": anime.Name}
		update := bson.M{
			"$set": bson.M{
				"name":        anime.Name,
				"synopsis":    anime.Synopsis,
				"genre":       anime.Genre,
				"score":       anime.Score,
				"imageUrl":    anime.ImageUrl,
				"bannerUrl":   anime.BannerUrl,
				"year":        anime.Year,
				"type":        anime.Type,
				"information": anime.Information,
				"updatedAt":   time.Now(),
			},
			"$setOnInsert": bson.M{
				"createdAt": time.Now(),
			},
		}

		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Error saving anime %s: %v", anime.Name, err)
		}
	}
	log.Printf("✅ Saved %d anime to database", len(animes))
}

// upgradeAnimeImages upgrades low-quality images in background
func upgradeAnimeImages(animes []models.Anime) {
	for _, anime := range animes {
		if isLowQualityImage(anime.ImageUrl) {
			go upgradeAnimeImage(anime.Name)
		}
	}
}

func isLowQualityImage(imageUrl string) bool {
	if imageUrl == "" {
		return true
	}
	// Check for low quality indicators
	return len(imageUrl) < 50 || 
		   !contains(imageUrl, "extraLarge") && !contains(imageUrl, "large")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

func upgradeAnimeImage(animeName string) {
	// This would call the existing image upgrade service
	log.Printf("🔄 Upgrading image for: %s", animeName)
}