package services

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	model "github.com/Flack74/mongoapi/models"
)

func SmartSearch(search, genre, year, season, format, status, userID string) []primitive.M {
	// First, try local search
	localResults := FilterAnimes(search, genre, year, season, format, status, userID)
	
	// If we have results or no search term, return local results
	if len(localResults) > 0 || search == "" {
		return localResults
	}
	
	// If no local results and we have a search term, try external search
	log.Printf("No local results for search '%s', trying external APIs", search)
	
	// Try to import from external APIs
	if search != "" {
		// Try Jikan API first
		if count, err := ImportSearchResults(search); err == nil && count > 0 {
			log.Printf("Imported %d anime from external search", count)
			// Search again after import
			return FilterAnimes(search, genre, year, season, format, status, userID)
		}
		
		// Try AniList as fallback
		if err := ImportFromAniList(search); err == nil {
			log.Printf("Imported anime from AniList for search '%s'", search)
			// Search again after import
			return FilterAnimes(search, genre, year, season, format, status, userID)
		}
	}
	
	// Return empty results if nothing found
	return []primitive.M{}
}

func ImportSearchResults(searchTerm string) (int, error) {
	// Use existing Jikan import but with search
	resp, err := http.Get("https://api.jikan.moe/v4/anime?q=" + searchTerm + "&limit=10")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var jikanResp JikanResponse
	if err := json.NewDecoder(resp.Body).Decode(&jikanResp); err != nil {
		return 0, err
	}

	return importAnimeList(jikanResp.Data)
}

func ImportFromAniList(searchTerm string) error {
	// Get anime data from AniList
	anilistData, err := GetAnimeFromAniList(searchTerm)
	if err != nil {
		return err
	}
	
	// Check if already exists
	existing, _ := FindAnimeByName(anilistData.Data.Media.Title.Romaji)
	if existing != nil {
		return nil // Already exists
	}
	
	// Create anime from AniList data
	anime := model.Anime{
		Name:      anilistData.Data.Media.Title.Romaji,
		Type:      model.AnimeType(anilistData.Data.Media.Format),
		Score:     0, // User will rate later
		Status:    "plan-to-watch",
		Genre:     anilistData.Data.Media.Genres,
		Notes:     truncateString(anilistData.Data.Media.Description, 500),
		Year:      anilistData.Data.Media.StartDate.Year,
		Season:    model.Season(anilistData.Data.Media.Season),
		ImageUrl:  anilistData.Data.Media.CoverImage.Large,
		BannerUrl: anilistData.Data.Media.BannerImage,
		Progress: model.Progress{
			Watched: 0,
			Total:   0,
		},
	}
	
	return InsertOneAnime(anime)
}

func SearchJikanAPI(query string) ([]JikanAnime, error) {
	resp, err := http.Get("https://api.jikan.moe/v4/anime?q=" + query + "&limit=10")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jikanResp JikanResponse
	if err := json.NewDecoder(resp.Body).Decode(&jikanResp); err != nil {
		return nil, err
	}

	return jikanResp.Data, nil
}