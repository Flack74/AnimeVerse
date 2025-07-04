package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"animeverse/config"
	"animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AniListAnimeResponse struct {
	Data struct {
		Page struct {
			Media []struct {
				ID          int    `json:"id"`
				Title       struct {
					Romaji  string `json:"romaji"`
					English string `json:"english"`
				} `json:"title"`
				Description string   `json:"description"`
				Genres      []string `json:"genres"`
				Episodes    int      `json:"episodes"`
				Status      string   `json:"status"`
				Format      string   `json:"format"`
				StartDate   struct {
					Year int `json:"year"`
				} `json:"startDate"`
				Season      string `json:"season"`
				AverageScore int   `json:"averageScore"`
				CoverImage  struct {
					Large string `json:"large"`
				} `json:"coverImage"`
				BannerImage string `json:"bannerImage"`
			} `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

func EnhanceAnimeFromAPI(animeName string) (*models.Anime, error) {
	// First try AniList
	anime, err := fetchFromAniList(animeName)
	if err == nil && anime != nil {
		// Save to database
		savedAnime, err := saveAnimeToDatabase(anime)
		if err == nil {
			// Cache the anime
			SetCachedAnime(savedAnime.ID.Hex(), savedAnime)
			return savedAnime, nil
		}
	}
	
	return nil, err
}

func fetchFromAniList(animeName string) (*models.Anime, error) {
	query := `
	query ($search: String) {
		Page(page: 1, perPage: 1) {
			media(search: $search, type: ANIME) {
				id
				title {
					romaji
					english
				}
				description(asHtml: false)
				genres
				episodes
				status
				format
				startDate {
					year
				}
				season
				averageScore
				coverImage {
					large
				}
				bannerImage
			}
		}
	}`

	variables := map[string]interface{}{
		"search": animeName,
	}

	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var anilistResp AniListAnimeResponse
	err = json.NewDecoder(resp.Body).Decode(&anilistResp)
	if err != nil {
		return nil, err
	}

	if len(anilistResp.Data.Page.Media) == 0 {
		return nil, nil
	}

	media := anilistResp.Data.Page.Media[0]
	
	// Convert to our anime model
	anime := &models.Anime{
		Name:      getPreferredTitle(media.Title.English, media.Title.Romaji),
		Type:      convertFormat(media.Format),
		Score:     float64(media.AverageScore) / 10.0,
		Genre:     media.Genres,
		Notes:     cleanDescription(media.Description),
		ImageUrl:  media.CoverImage.Large,
		BannerUrl: media.BannerImage,
		AniListID: media.ID,
		Year:      media.StartDate.Year,
		Season:    convertSeason(media.Season),
		Status:    convertStatus(media.Status),
		Progress: models.Progress{
			Total: media.Episodes,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return anime, nil
}

func saveAnimeToDatabase(anime *models.Anime) (*models.Anime, error) {
	anime.ID = primitive.NewObjectID()
	
	_, err := config.Collection.InsertOne(context.Background(), anime)
	if err != nil {
		return nil, err
	}
	
	return anime, nil
}

func getPreferredTitle(english, romaji string) string {
	if english != "" {
		return english
	}
	return romaji
}

func convertFormat(format string) models.AnimeType {
	switch strings.ToUpper(format) {
	case "TV":
		return models.SeriesType
	case "MOVIE":
		return models.MovieType
	case "ONA":
		return models.ONAType
	default:
		return models.SeriesType
	}
}

func convertScore(score int) int {
	if score > 0 {
		return score / 10 // Convert from 100 scale to 10 scale
	}
	return 0
}

func convertSeason(season string) models.Season {
	switch strings.ToUpper(season) {
	case "WINTER":
		return models.Winter
	case "SPRING":
		return models.Spring
	case "SUMMER":
		return models.Summer
	case "FALL":
		return models.Fall
	default:
		return models.Spring
	}
}

func convertStatus(status string) models.WatchStatus {
	switch strings.ToUpper(status) {
	case "FINISHED":
		return models.Completed
	case "RELEASING":
		return models.Watching
	default:
		return models.PlanToWatch
	}
}

func cleanDescription(description string) string {
	if description == "" {
		return "No description available."
	}
	
	// Remove HTML tags and limit length
	cleaned := strings.ReplaceAll(description, "<br>", " ")
	cleaned = strings.ReplaceAll(cleaned, "<i>", "")
	cleaned = strings.ReplaceAll(cleaned, "</i>", "")
	
	if len(cleaned) > 500 {
		cleaned = cleaned[:500] + "..."
	}
	
	return cleaned
}

func UpdateAnimeWithAPIData(animeID string, animeName string) error {
	// Fetch from API
	apiAnime, err := fetchFromAniList(animeName)
	if err != nil || apiAnime == nil {
		return err
	}
	
	// Update existing anime in database
	objectID, err := primitive.ObjectIDFromHex(animeID)
	if err != nil {
		return err
	}
	
	update := bson.M{
		"$set": bson.M{
			"notes":     apiAnime.Notes,
			"genre":     apiAnime.Genre,
			"imageUrl":  apiAnime.ImageUrl,
			"bannerUrl": apiAnime.BannerUrl,
			"year":      apiAnime.Year,
			"season":    apiAnime.Season,
			"anilist_id": apiAnime.AniListID,
			"progress.total": apiAnime.Progress.Total,
			"updated_at": time.Now(),
		},
	}
	
	_, err = config.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)
	
	if err == nil {
		// Invalidate cache
		InvalidateAnimeCache(animeID)
	}
	
	return err
}