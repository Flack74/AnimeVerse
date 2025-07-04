package services

import (
	"context"
	"time"

	"animeverse/cache"
	"animeverse/config"
	"animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAnimeWithFallback tries database first, then AniList if needed
func GetAnimeWithFallback(animeName string) (*models.Anime, error) {
	// 1. Check cache first
	cacheKey := "anime_fallback:" + animeName
	var cached models.Anime
	if err := cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// 2. Try database first
	anime, err := findAnimeInDatabase(animeName)
	if err == nil && isAnimeComplete(anime) {
		// Cache complete anime for 1 hour
		cache.Set(cacheKey, *anime, time.Hour)
		return anime, nil
	}

	// 3. Enhance from AniList if incomplete or not found
	if anime != nil {
		// Enhance existing anime
		enhanced, err := enhanceExistingAnime(anime, animeName)
		if err == nil {
			cache.Set(cacheKey, *enhanced, time.Hour)
			return enhanced, nil
		}
		return anime, nil // Return basic anime if enhancement fails
	}

	// 4. Create new anime from AniList
	newAnime, err := createAnimeFromAniList(animeName)
	if err != nil {
		return nil, err
	}

	// Cache new anime for 1 hour
	cache.Set(cacheKey, *newAnime, time.Hour)
	return newAnime, nil
}

func findAnimeInDatabase(animeName string) (*models.Anime, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": animeName, "$options": "i"}},
			{"alternative_titles.english": bson.M{"$regex": animeName, "$options": "i"}},
			{"alternative_titles.synonyms": bson.M{"$in": []string{animeName}}},
		},
	}

	var anime models.Anime
	err := config.Collection.FindOne(context.Background(), filter).Decode(&anime)
	if err != nil {
		return nil, err
	}

	return &anime, nil
}

func isAnimeComplete(anime *models.Anime) bool {
	// Check if anime has essential data
	return anime.Synopsis != "" &&
		len(anime.Genre) > 0 &&
		anime.ImageUrl != "" &&
		anime.Information.Episodes > 0 &&
		anime.Information.Status != "" &&
		len(anime.Characters) > 0
}

func enhanceExistingAnime(anime *models.Anime, animeName string) (*models.Anime, error) {
	// Get enhanced data from AniList
	enhanced, err := fetchFullAnimeDataFromAniList(animeName)
	if err != nil {
		return nil, err
	}

	// Merge with existing data (keep user data, enhance missing fields)
	if anime.Synopsis == "" && enhanced.Information.Status != "" {
		anime.Synopsis = enhanced.Information.Status // Use description from AniList
	}
	
	if len(anime.Genre) == 0 {
		// Get genres from AniList response
		anime.Genre = []string{"Action", "Drama"} // Placeholder - should come from enhanced data
	}

	if anime.ImageUrl == "" {
		// Get high-quality image from AniList
		anime.ImageUrl = "https://via.placeholder.com/300x400" // Placeholder
	}

	if anime.BannerUrl == "" {
		// Get high-quality banner from AniList
		anime.BannerUrl = "https://via.placeholder.com/800x300" // Placeholder
	}

	// Update missing detailed information
	if anime.Information.Episodes == 0 {
		anime.Information = enhanced.Information
	}
	
	if len(anime.Characters) == 0 {
		anime.Characters = enhanced.Characters
	}
	
	if len(anime.Staff) == 0 {
		anime.Staff = enhanced.Staff
	}

	anime.AlternativeTitles = enhanced.AlternativeTitles
	anime.Statistics = enhanced.Statistics
	anime.Related = enhanced.Related

	// Update in database
	updateAnimeInDatabase(anime)

	return anime, nil
}

func createAnimeFromAniList(animeName string) (*models.Anime, error) {
	// Fetch complete data from AniList
	anilistData, err := fetchFromAniList(animeName)
	if err != nil {
		return nil, err
	}

	// Create new anime with complete data
	anime := &models.Anime{
		ID:        primitive.NewObjectID(),
		Name:      anilistData.Name,
		Synopsis:  anilistData.Notes,
		Genre:     anilistData.Genre,
		Score:     anilistData.Score,
		ImageUrl:  anilistData.ImageUrl,
		BannerUrl: anilistData.BannerUrl,
		Year:      anilistData.Year,
		Season:    anilistData.Season,
		Type:      anilistData.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Get enhanced details
	enhanced, err := fetchFullAnimeDataFromAniList(animeName)
	if err == nil {
		anime.AlternativeTitles = enhanced.AlternativeTitles
		anime.Information = enhanced.Information
		anime.Statistics = enhanced.Statistics
		anime.Characters = enhanced.Characters
		anime.Staff = enhanced.Staff
		anime.Related = enhanced.Related
	}

	// Save to database
	_, err = config.Collection.InsertOne(context.Background(), anime)
	if err != nil {
		return nil, err
	}

	return anime, nil
}

func updateAnimeInDatabase(anime *models.Anime) error {
	anime.UpdatedAt = time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"synopsis":           anime.Synopsis,
			"genre":              anime.Genre,
			"imageUrl":           anime.ImageUrl,
			"bannerUrl":          anime.BannerUrl,
			"alternative_titles": anime.AlternativeTitles,
			"information":        anime.Information,
			"statistics":         anime.Statistics,
			"characters":         anime.Characters,
			"staff":              anime.Staff,
			"related":            anime.Related,
			"updated_at":         anime.UpdatedAt,
		},
	}

	_, err := config.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": anime.ID},
		update,
	)

	return err
}

// GetSpotlightWithDatabase gets spotlight anime from database first
func GetSpotlightWithDatabase() ([]SpotlightAnime, error) {
	// Check cache first
	var cached []SpotlightAnime
	if err := cache.Get(SPOTLIGHT_CACHE_KEY, &cached); err == nil {
		// Upgrade images if needed
		for i := range cached {
			if needsImageUpgrade(cached[i].BannerUrl, cached[i].ImageUrl) {
				newBanner, newCover, err := UpgradeImageQualityAndSave(cached[i].ID, cached[i].Name, cached[i].BannerUrl, cached[i].ImageUrl)
				if err == nil {
					cached[i].BannerUrl = newBanner
					cached[i].ImageUrl = newCover
				}
			}
		}
		return cached, nil
	}

	// Try database first
	dbSpotlight, err := getSpotlightFromDatabase()
	if err == nil && len(dbSpotlight) >= 10 {
		// Upgrade images for database results
		for i := range dbSpotlight {
			if needsImageUpgrade(dbSpotlight[i].BannerUrl, dbSpotlight[i].ImageUrl) {
				newBanner, newCover, err := UpgradeImageQualityAndSave(dbSpotlight[i].ID, dbSpotlight[i].Name, dbSpotlight[i].BannerUrl, dbSpotlight[i].ImageUrl)
				if err == nil {
					dbSpotlight[i].BannerUrl = newBanner
					dbSpotlight[i].ImageUrl = newCover
				}
			}
		}
		cache.Set(SPOTLIGHT_CACHE_KEY, dbSpotlight, time.Hour)
		return dbSpotlight, nil
	}

	// Fallback to AniList and save to database
	anilistSpotlight, err := fetchSpotlightFromAniList()
	if err != nil {
		return nil, err
	}

	// Save to database for future use
	saveSpotlightToDatabase(anilistSpotlight)

	cache.Set(SPOTLIGHT_CACHE_KEY, anilistSpotlight, time.Hour)
	return anilistSpotlight, nil
}

func getSpotlightFromDatabase() ([]SpotlightAnime, error) {
	filter := bson.M{
		"score": bson.M{"$gte": 8.5},
		"imageUrl": bson.M{"$ne": ""},
		"bannerUrl": bson.M{"$ne": ""},
	}

	cursor, err := config.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var spotlight []SpotlightAnime
	for cursor.Next(context.Background()) {
		var anime models.Anime
		if err := cursor.Decode(&anime); err != nil {
			continue
		}

		spotlight = append(spotlight, SpotlightAnime{
			ID:          anime.ID.Hex(),
			Name:        anime.Name,
			Score:       anime.Score,
			Year:        anime.Year,
			ImageUrl:    anime.ImageUrl,
			BannerUrl:   anime.BannerUrl,
			Description: anime.Synopsis,
			Genres:      anime.Genre,
		})

		if len(spotlight) >= 10 {
			break
		}
	}

	return spotlight, nil
}

func saveSpotlightToDatabase(spotlight []SpotlightAnime) {
	for _, item := range spotlight {
		// Check if already exists
		filter := bson.M{"name": item.Name}
		var existing models.Anime
		err := config.Collection.FindOne(context.Background(), filter).Decode(&existing)
		
		if err != nil {
			// Create new anime
			anime := models.Anime{
				ID:        primitive.NewObjectID(),
				Name:      item.Name,
				Score:     item.Score,
				Year:      item.Year,
				ImageUrl:  item.ImageUrl,
				BannerUrl: item.BannerUrl,
				Synopsis:  item.Description,
				Genre:     item.Genres,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			config.Collection.InsertOne(context.Background(), anime)
		} else {
			// Update existing with high-quality images
			update := bson.M{
				"$set": bson.M{
					"imageUrl":   item.ImageUrl,
					"bannerUrl":  item.BannerUrl,
					"score":      item.Score,
					"updated_at": time.Now(),
				},
			}
			config.Collection.UpdateOne(context.Background(), filter, update)
		}
	}
}