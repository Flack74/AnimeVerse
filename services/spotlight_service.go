package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"


	"animeverse/config"
	"animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const SPOTLIGHT_CACHE_KEY = "spotlight:anime"

type SpotlightAnime struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	Year        int     `json:"year"`
	ImageUrl    string  `json:"imageUrl"`
	BannerUrl   string  `json:"bannerUrl"`
	Description string  `json:"description"`
	Genres      []string `json:"genres"`
}

func GetSpotlightAnime() ([]SpotlightAnime, error) {
	return GetSpotlightWithDatabase()
}

func fetchSpotlightFromAniList() ([]SpotlightAnime, error) {
	query := `{
		Page(page: 1, perPage: 10) {
			media(type: ANIME, sort: [SCORE_DESC, POPULARITY_DESC], averageScore_greater: 85) {
				id
				title { romaji english }
				averageScore
				startDate { year }
				coverImage { extraLarge }
				bannerImage
				description(asHtml: false)
				genres
			}
		}
	}`

	requestBody := map[string]interface{}{
		"query": query,
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

	var anilistResp struct {
		Data struct {
			Page struct {
				Media []struct {
					ID          int    `json:"id"`
					Title       struct {
						Romaji  string `json:"romaji"`
						English string `json:"english"`
					} `json:"title"`
					AverageScore int      `json:"averageScore"`
					StartDate    struct {
						Year int `json:"year"`
					} `json:"startDate"`
					CoverImage struct {
						ExtraLarge string `json:"extraLarge"`
					} `json:"coverImage"`
					BannerImage string   `json:"bannerImage"`
					Description string   `json:"description"`
					Genres      []string `json:"genres"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&anilistResp); err != nil {
		return nil, err
	}

	var spotlight []SpotlightAnime
	for _, media := range anilistResp.Data.Page.Media {
		title := media.Title.English
		if title == "" {
			title = media.Title.Romaji
		}

		spotlight = append(spotlight, SpotlightAnime{
			ID:          fmt.Sprintf("anilist_%d", media.ID),
			Name:        title,
			Score:       float64(media.AverageScore) / 10.0,
			Year:        media.StartDate.Year,
			ImageUrl:    media.CoverImage.ExtraLarge,
			BannerUrl:   media.BannerImage,
			Description: cleanDescription(media.Description),
			Genres:      media.Genres,
		})
	}

	return spotlight, nil
}

func GetTopRatedMixed() ([]models.Anime, error) {
	// Get top rated from different genres
	genres := []string{"Action", "Drama", "Comedy", "Romance", "Fantasy", "Thriller"}
	var allAnime []models.Anime

	for _, genre := range genres {
		filter := bson.M{
			"genre": bson.M{"$in": []string{genre}},
			"score": bson.M{"$gte": 7.0},
		}
		
		opts := options.Find().SetLimit(2).SetSort(bson.D{{"score", -1}})
		
		cursor, err := config.Collection.Find(context.Background(), filter, opts)
		if err != nil {
			continue
		}
		
		var genreAnime []models.Anime
		if err := cursor.All(context.Background(), &genreAnime); err == nil {
			allAnime = append(allAnime, genreAnime...)
		}
		cursor.Close(context.Background())
	}

	// Verify scores with external API
	return verifyAndEnhanceScores(allAnime)
}

func verifyAndEnhanceScores(animes []models.Anime) ([]models.Anime, error) {
	var verified []models.Anime
	
	for _, anime := range animes {
		// Skip if already has high confidence score
		if anime.Score >= 8.0 {
			verified = append(verified, anime)
			continue
		}

		// Verify with AniList
		externalScore, err := getScoreFromAniList(anime.Name)
		if err != nil {
			continue
		}

		// Update if external score is significantly different
		if externalScore > 0 && abs(anime.Score, externalScore) > 1.0 {
			anime.Score = externalScore
			updateAnimeScore(anime.ID, externalScore)
		}

		if anime.Score >= 7.0 {
			verified = append(verified, anime)
		}
	}

	return verified, nil
}

func getScoreFromAniList(animeName string) (float64, error) {
	query := fmt.Sprintf(`{
		Page(page: 1, perPage: 1) {
			media(search: "%s", type: ANIME) {
				averageScore
			}
		}
	}`, animeName)

	requestBody := map[string]interface{}{
		"query": query,
	}

	jsonData, _ := json.Marshal(requestBody)
	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Page struct {
				Media []struct {
					AverageScore int `json:"averageScore"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if len(result.Data.Page.Media) > 0 && result.Data.Page.Media[0].AverageScore > 0 {
		return float64(result.Data.Page.Media[0].AverageScore) / 10.0, nil
	}

	return 0, fmt.Errorf("no score found")
}

func updateAnimeScore(animeID primitive.ObjectID, newScore float64) {
	update := bson.M{
		"$set": bson.M{
			"score":      newScore,
			"updated_at": time.Now(),
		},
	}
	
	config.Collection.UpdateOne(context.Background(), bson.M{"_id": animeID}, update)
}

func abs(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}