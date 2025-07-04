package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"animeverse/config"
	"go.mongodb.org/mongo-driver/bson"
)

type AniListResponse struct {
	Data struct {
		Media struct {
			ID          int    `json:"id"`
			Title       struct {
				Romaji string `json:"romaji"`
				English string `json:"english"`
			} `json:"title"`
			StartDate struct {
				Year int `json:"year"`
			} `json:"startDate"`
			Season      string `json:"season"`
			Format      string `json:"format"`
			Genres      []string `json:"genres"`
			Description string `json:"description"`
			CoverImage  struct {
				Large string `json:"large"`
			} `json:"coverImage"`
			BannerImage string `json:"bannerImage"`
		} `json:"Media"`
	} `json:"data"`
}

func GetAnimeFromAniList(animeName string) (*AniListResponse, error) {
	query := `
	query ($search: String) {
		Media (search: $search, type: ANIME) {
			id
			title {
				romaji
				english
			}
			startDate {
				year
			}
			season
			format
			genres
			description
			coverImage {
				large
			}
			bannerImage
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

	var anilistResp AniListResponse
	err = json.NewDecoder(resp.Body).Decode(&anilistResp)
	if err != nil {
		return nil, err
	}

	return &anilistResp, nil
}

func BackfillAnimeData(animeName string) error {
	// Get anime from AniList
	anilistData, err := GetAnimeFromAniList(animeName)
	if err != nil {
		return err
	}

	// Update MongoDB record with missing data
	filter := bson.M{"name": bson.M{"$regex": "^" + animeName + "$", "$options": "i"}}
	
	updateFields := bson.M{}
	
	if anilistData.Data.Media.StartDate.Year > 0 {
		updateFields["year"] = anilistData.Data.Media.StartDate.Year
	}
	
	if anilistData.Data.Media.Season != "" {
		season := strings.Title(strings.ToLower(anilistData.Data.Media.Season))
		updateFields["season"] = season
	}
	
	if anilistData.Data.Media.CoverImage.Large != "" {
		updateFields["imageUrl"] = anilistData.Data.Media.CoverImage.Large
	}
	
	if anilistData.Data.Media.BannerImage != "" {
		updateFields["bannerUrl"] = anilistData.Data.Media.BannerImage
	}
	
	if len(anilistData.Data.Media.Genres) > 0 {
		updateFields["genre"] = anilistData.Data.Media.Genres
	}
	
	update := bson.M{"$set": updateFields}

	_, err = config.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func BackfillAllMissingData() (int, error) {
	// Find anime with missing year or season data
	filter := bson.M{
		"$or": []bson.M{
			{"year": bson.M{"$exists": false}},
			{"year": 0},
			{"season": bson.M{"$exists": false}},
			{"season": ""},
		},
	}

	cur, err := config.Collection.Find(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	defer cur.Close(context.Background())

	count := 0
	for cur.Next(context.Background()) {
		var anime bson.M
		if err := cur.Decode(&anime); err != nil {
			continue
		}

		if name, ok := anime["name"].(string); ok {
			if err := BackfillAnimeData(name); err == nil {
				count++
			}
			// Rate limiting - AniList allows 90 requests per minute
			time.Sleep(700 * time.Millisecond)
		}
	}

	return count, nil
}