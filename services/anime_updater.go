package services

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"animeverse/config"
	model "animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AniListCurrentSeason represents current season anime from AniList
type AniListCurrentSeason struct {
	Data struct {
		Page struct {
			Media []struct {
				ID          int    `json:"id"`
				Title       struct {
					Romaji string `json:"romaji"`
					English string `json:"english"`
				} `json:"title"`
				Format      string   `json:"format"`
				Status      string   `json:"status"`
				Episodes    int      `json:"episodes"`
				Season      string   `json:"season"`
				SeasonYear  int      `json:"seasonYear"`
				AverageScore int     `json:"averageScore"`
				Genres      []string `json:"genres"`
				CoverImage  struct {
					Large string `json:"large"`
				} `json:"coverImage"`
				BannerImage string `json:"bannerImage"`
			} `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

// UpdateCurrentSeasonAnime fetches and updates current season anime
func UpdateCurrentSeasonAnime() (int, error) {
	log.Println("Updating current season anime from AniList...")
	
	query := `
	query ($season: MediaSeason, $year: Int, $page: Int) {
		Page(page: $page, perPage: 50) {
			media(season: $season, seasonYear: $year, type: ANIME, sort: POPULARITY_DESC) {
				id
				title {
					romaji
					english
				}
				format
				status
				episodes
				season
				seasonYear
				averageScore
				genres
				coverImage {
					large
				}
				bannerImage
			}
		}
	}`

	currentYear := time.Now().Year()
	currentSeason := getCurrentSeason()
	
	variables := map[string]interface{}{
		"season": currentSeason,
		"year":   currentYear,
		"page":   1,
	}

	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var anilistResp AniListCurrentSeason
	if err := json.NewDecoder(resp.Body).Decode(&anilistResp); err != nil {
		return 0, err
	}

	updated := 0
	for _, media := range anilistResp.Data.Page.Media {
		if media.Title.Romaji == "" {
			continue
		}

		// Check if anime already exists
		title := media.Title.Romaji
		if media.Title.English != "" {
			title = media.Title.English
		}

		filter := bson.M{
			"$or": []bson.M{
				{"name": title},
				{"anilist_id": media.ID},
			},
		}

		var existingAnime model.Anime
		err := config.Collection.FindOne(context.Background(), filter).Decode(&existingAnime)
		
		if err == mongo.ErrNoDocuments {
			// Create new anime
			newAnime := model.Anime{
				Name:      title,
				Type:      convertAniListFormat(media.Format),
				Score:     float64(media.AverageScore) / 10.0, // Convert from 100 scale to 10 scale
				Progress:  model.Progress{Total: media.Episodes},
				Status:    convertAniListStatus(media.Status),
				Genre:     media.Genres,
				ImageUrl:  media.CoverImage.Large,
				BannerUrl: media.BannerImage,
				AniListID: media.ID,
				Year:      media.SeasonYear,
				Season:    model.Season(media.Season),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			_, err := config.Collection.InsertOne(context.Background(), newAnime)
			if err != nil {
				log.Printf("Error inserting new anime %s: %v", title, err)
				continue
			}
			updated++
		} else if err == nil {
			// Update existing anime with new data
			update := bson.M{
				"$set": bson.M{
					"score":      float64(media.AverageScore) / 10.0,
					"status":     convertAniListStatus(media.Status),
					"imageUrl":   media.CoverImage.Large,
					"bannerUrl":  media.BannerImage,
					"anilist_id": media.ID,
					"updated_at": time.Now(),
				},
			}

			_, err := config.Collection.UpdateOne(context.Background(), bson.M{"_id": existingAnime.ID}, update)
			if err != nil {
				log.Printf("Error updating anime %s: %v", title, err)
				continue
			}
			updated++
		}
	}

	log.Printf("Updated %d current season anime", updated)
	return updated, nil
}

func getCurrentSeason() string {
	month := time.Now().Month()
	switch {
	case month >= 12 || month <= 2:
		return "WINTER"
	case month >= 3 && month <= 5:
		return "SPRING"
	case month >= 6 && month <= 8:
		return "SUMMER"
	case month >= 9 && month <= 11:
		return "FALL"
	default:
		return "WINTER"
	}
}

func convertAniListFormat(format string) model.AnimeType {
	switch format {
	case "TV":
		return model.SeriesType
	case "MOVIE":
		return model.MovieType
	case "ONA", "OVA", "SPECIAL":
		return model.ONAType
	default:
		return model.SeriesType
	}
}

func convertAniListStatus(status string) model.WatchStatus {
	switch status {
	case "FINISHED":
		return model.Completed
	case "RELEASING":
		return model.Watching
	case "NOT_YET_RELEASED":
		return model.PlanToWatch
	default:
		return model.PlanToWatch
	}
}