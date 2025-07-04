package services

import (
	"context"
	"fmt"
	"time"

	"animeverse/cache"
	"animeverse/config"
	"animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAnimesWithRedisCache(limit int, offset int) ([]models.Anime, error) {
	cacheKey := fmt.Sprintf("animes:limit_%d:offset_%d", limit, offset)
	
	// Try Redis first
	var cached []models.Anime
	if err := cache.Get(cacheKey, &cached); err == nil {
		return cached, nil
	}

	// Get from database
	cursor, err := config.Collection.Find(context.Background(), bson.M{}, &options.FindOptions{
		Limit: &[]int64{int64(limit)}[0],
		Skip:  &[]int64{int64(offset)}[0],
		Sort:  bson.D{{"score", -1}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var animes []models.Anime
	if err := cursor.All(context.Background(), &animes); err != nil {
		return nil, err
	}

	// Cache for 10 minutes
	cache.Set(cacheKey, animes, 10*time.Minute)
	
	return animes, nil
}

func GetTrendingWithRedisCache() ([]models.Anime, error) {
	cacheKey := "trending_animes_fast"
	
	// Try Redis first
	var cached []models.Anime
	if err := cache.Get(cacheKey, &cached); err == nil {
		return cached, nil
	}

	// Get top rated from database
	cursor, err := config.Collection.Find(context.Background(), bson.M{
		"score": bson.M{"$gte": 7.0},
	}, &options.FindOptions{
		Limit: &[]int64{50}[0],
		Sort:  bson.D{{"score", -1}, {"year", -1}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var animes []models.Anime
	if err := cursor.All(context.Background(), &animes); err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	cache.Set(cacheKey, animes, 5*time.Minute)
	
	return animes, nil
}