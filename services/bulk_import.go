package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"animeverse/config"
	model "animeverse/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnimeOfflineDB struct {
	Sources      []string    `json:"sources"`
	Title        string      `json:"title"`
	Type         string      `json:"type"`
	Episodes     int         `json:"episodes"`
	Status       string      `json:"status"`
	AnimeSeason  interface{} `json:"animeSeason"`
	Year         int         `json:"year"`
	Tags         []string    `json:"tags"`
	Synonyms     []string    `json:"synonyms"`
	Relations    []string    `json:"relations"`
	Picture      string      `json:"picture"`
	Thumbnail    string      `json:"thumbnail"`
}

func BulkImportAnimeDatabase() (int, error) {
	log.Println("Starting bulk import from anime-offline-database...")
	
	// Fetch data from GitHub
	resp, err := http.Get("https://raw.githubusercontent.com/manami-project/anime-offline-database/refs/heads/master/anime-offline-database-minified.json")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch anime database: %v", err)
	}
	defer resp.Body.Close()

	var animeData struct {
		Data []AnimeOfflineDB `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&animeData); err != nil {
		return 0, fmt.Errorf("failed to decode JSON: %v", err)
	}

	log.Printf("Fetched %d anime from database", len(animeData.Data))

	// Convert to our model and batch insert
	var animes []interface{}
	for _, item := range animeData.Data {
		if item.Title == "" {
			continue
		}

		// Extract MAL and AniList IDs from sources
		malID, anilistID := extractIDs(item.Sources)
		
		// Convert type
		animeType := convertType(item.Type)
		
		// Convert status
		status := convertBulkStatus(item.Status)
		
		// Generate score based on various factors
		score := generateScore(item)
		
		// Extract season from interface
		season := ""
		if item.AnimeSeason != nil {
			if seasonMap, ok := item.AnimeSeason.(map[string]interface{}); ok {
				if s, ok := seasonMap["season"].(string); ok {
					season = s
				}
			}
		}

		anime := model.Anime{
			Name:      item.Title,
			Type:      animeType,
			Score:     float64(score),
			Progress:  model.Progress{Total: item.Episodes},
			Status:    status,
			Genre:     item.Tags,
			Notes:     strings.Join(item.Synonyms, ", "),
			ImageUrl:  item.Picture,
			BannerUrl: item.Thumbnail,
			MALID:     malID,
			AniListID: anilistID,
			Year:      item.Year,
			Season:    model.Season(season),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		animes = append(animes, anime)

		// Batch insert every 1000 records
		if len(animes) >= 1000 {
			if err := insertBatch(animes); err != nil {
				log.Printf("Error inserting batch: %v", err)
			}
			animes = animes[:0] // Clear slice
		}
	}

	// Insert remaining records
	if len(animes) > 0 {
		if err := insertBatch(animes); err != nil {
			log.Printf("Error inserting final batch: %v", err)
		}
	}

	log.Printf("Bulk import completed. Imported %d anime", len(animeData.Data))
	return len(animeData.Data), nil
}

func extractIDs(sources []string) (int, int) {
	var malID, anilistID int
	
	for _, source := range sources {
		if strings.Contains(source, "myanimelist.net") {
			parts := strings.Split(source, "/")
			for i, part := range parts {
				if part == "anime" && i+1 < len(parts) {
					if id, err := strconv.Atoi(parts[i+1]); err == nil {
						malID = id
					}
					break
				}
			}
		} else if strings.Contains(source, "anilist.co") {
			parts := strings.Split(source, "/")
			for i, part := range parts {
				if part == "anime" && i+1 < len(parts) {
					if id, err := strconv.Atoi(parts[i+1]); err == nil {
						anilistID = id
					}
					break
				}
			}
		}
	}
	
	return malID, anilistID
}

func convertType(t string) model.AnimeType {
	switch strings.ToUpper(t) {
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

func convertBulkStatus(s string) model.WatchStatus {
	switch strings.ToLower(s) {
	case "finished", "completed":
		return model.Completed
	case "ongoing", "currently airing":
		return model.Watching
	default:
		return model.PlanToWatch
	}
}

func generateScore(item AnimeOfflineDB) int {
	score := 5 // Base score
	
	// Boost score for recent anime
	if item.Year >= 2020 {
		score += 2
	} else if item.Year >= 2015 {
		score += 1
	}
	
	// Boost for popular tags
	popularTags := []string{"Action", "Adventure", "Comedy", "Drama", "Fantasy", "Romance", "Shounen"}
	for _, tag := range item.Tags {
		for _, popular := range popularTags {
			if strings.EqualFold(tag, popular) {
				score += 1
				break
			}
		}
	}
	
	// Cap at 10
	if score > 10 {
		score = 10
	}
	
	return score
}

func insertBatch(animes []interface{}) error {
	if len(animes) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use ordered=false for better performance
	opts := &options.InsertManyOptions{}
	opts.SetOrdered(false)

	_, err := config.Collection.InsertMany(ctx, animes, opts)
	if err != nil {
		// Log but don't fail on duplicate key errors
		if mongo.IsDuplicateKeyError(err) {
			log.Printf("Some duplicates found in batch, continuing...")
			return nil
		}
		return err
	}

	log.Printf("Inserted batch of %d anime", len(animes))
	return nil
}