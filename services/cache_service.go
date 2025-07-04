package services

import (
	"fmt"
	"time"

	"animeverse/cache"
	"animeverse/models"
)

const (
	ANIME_CACHE_PREFIX    = "anime:"
	TRENDING_CACHE_KEY    = "trending:anime"
	MOVIES_CACHE_KEY      = "movies:anime"
	POPULAR_CACHE_KEY     = "popular:anime"
	SCHEDULE_CACHE_KEY    = "schedule:anime"
	CACHE_DURATION        = 5 * time.Minute
	LONG_CACHE_DURATION   = 30 * time.Minute
)

func GetCachedAnimes(key string) ([]models.Anime, error) {
	var animes []models.Anime
	err := cache.Get(key, &animes)
	return animes, err
}

func SetCachedAnimes(key string, animes []models.Anime, duration time.Duration) error {
	return cache.Set(key, animes, duration)
}

func GetCachedTrending() ([]models.Anime, error) {
	return GetCachedAnimes(TRENDING_CACHE_KEY)
}

func SetCachedTrending(animes []models.Anime) error {
	return SetCachedAnimes(TRENDING_CACHE_KEY, animes, CACHE_DURATION)
}

func GetCachedMovies() ([]models.Anime, error) {
	return GetCachedAnimes(MOVIES_CACHE_KEY)
}

func SetCachedMovies(animes []models.Anime) error {
	return SetCachedAnimes(MOVIES_CACHE_KEY, animes, LONG_CACHE_DURATION)
}

func GetCachedPopular() ([]models.Anime, error) {
	return GetCachedAnimes(POPULAR_CACHE_KEY)
}

func SetCachedPopular(animes []models.Anime) error {
	return SetCachedAnimes(POPULAR_CACHE_KEY, animes, CACHE_DURATION)
}

func GetCachedAnime(id string) (*models.Anime, error) {
	var anime models.Anime
	key := fmt.Sprintf("%s%s", ANIME_CACHE_PREFIX, id)
	err := cache.Get(key, &anime)
	if err != nil {
		return nil, err
	}
	return &anime, nil
}

func SetCachedAnime(id string, anime *models.Anime) error {
	key := fmt.Sprintf("%s%s", ANIME_CACHE_PREFIX, id)
	return cache.Set(key, anime, LONG_CACHE_DURATION)
}

func InvalidateAnimeCache(id string) error {
	key := fmt.Sprintf("%s%s", ANIME_CACHE_PREFIX, id)
	return cache.Delete(key)
}

func InvalidateAllCache() error {
	keys := []string{
		TRENDING_CACHE_KEY,
		MOVIES_CACHE_KEY,
		POPULAR_CACHE_KEY,
		SCHEDULE_CACHE_KEY,
	}
	
	for _, key := range keys {
		cache.Delete(key)
	}
	return nil
}