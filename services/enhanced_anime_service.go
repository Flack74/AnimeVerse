package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"animeverse/cache"
	"animeverse/config"
	"animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnhancedAnimeData struct {
	AlternativeTitles models.AlternativeTitles `json:"alternative_titles"`
	Information       models.AnimeInformation  `json:"information"`
	Statistics        models.AnimeStatistics   `json:"statistics"`
	Characters        []models.Character       `json:"characters"`
	Staff             []models.StaffMember     `json:"staff"`
	Themes            models.AnimeThemes       `json:"themes"`
	Related           []models.RelatedAnime    `json:"related"`
}

func EnhanceAnimeWithFullData(animeID string, animeName string) (*models.Anime, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("enhanced_anime:%s", animeID)
	var cached models.Anime
	if err := cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// Get base anime from database
	objectID, err := primitive.ObjectIDFromHex(animeID)
	if err != nil {
		return nil, err
	}

	var anime models.Anime
	err = config.Collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&anime)
	if err != nil {
		return nil, err
	}

	// Enhance with AniList data
	enhanced, err := fetchFullAnimeDataFromAniList(animeName)
	if err != nil {
		return &anime, nil // Return basic anime if enhancement fails
	}

	// Merge enhanced data
	anime.AlternativeTitles = enhanced.AlternativeTitles
	anime.Information = enhanced.Information
	anime.Statistics = enhanced.Statistics
	anime.Characters = enhanced.Characters
	anime.Staff = enhanced.Staff
	anime.Themes = enhanced.Themes
	anime.Related = enhanced.Related

	// Update in database
	updateData := bson.M{
		"$set": bson.M{
			"alternative_titles": enhanced.AlternativeTitles,
			"information":        enhanced.Information,
			"statistics":         enhanced.Statistics,
			"characters":         enhanced.Characters,
			"staff":              enhanced.Staff,
			"themes":             enhanced.Themes,
			"related":            enhanced.Related,
			"updated_at":         time.Now(),
		},
	}

	config.Collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, updateData)

	// Cache for 24 hours
	cache.Set(cacheKey, anime, 24*time.Hour)

	return &anime, nil
}

func fetchFullAnimeDataFromAniList(animeName string) (*EnhancedAnimeData, error) {
	query := fmt.Sprintf(`{
		Page(page: 1, perPage: 1) {
			media(search: "%s", type: ANIME) {
				title {
					romaji
					english
					native
				}
				synonyms
				episodes
				status
				startDate {
					year
					month
					day
				}
				endDate {
					year
					month
					day
				}
				season
				seasonYear
				format
				source
				duration
				averageScore
				popularity
				favourites
				rankings {
					rank
					type
					season
					year
				}
				studios {
					nodes {
						name
					}
				}
				genres
				characters(page: 1, perPage: 10, sort: ROLE) {
					edges {
						node {
							name {
								full
							}
							image {
								medium
							}
						}
						role
						voiceActors(language: JAPANESE) {
							name {
								full
							}
							image {
								medium
							}
						}
					}
				}
				staff(page: 1, perPage: 8) {
					edges {
						node {
							name {
								full
							}
							image {
								medium
							}
						}
						role
					}
				}
				relations {
					edges {
						node {
							title {
								romaji
								english
							}
							coverImage {
								medium
							}
						}
						relationType
					}
				}
			}
		}
	}`, animeName)

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
					Title struct {
						Romaji  string `json:"romaji"`
						English string `json:"english"`
						Native  string `json:"native"`
					} `json:"title"`
					Synonyms     []string `json:"synonyms"`
					Episodes     int      `json:"episodes"`
					Status       string   `json:"status"`
					StartDate    struct {
						Year  int `json:"year"`
						Month int `json:"month"`
						Day   int `json:"day"`
					} `json:"startDate"`
					EndDate struct {
						Year  int `json:"year"`
						Month int `json:"month"`
						Day   int `json:"day"`
					} `json:"endDate"`
					Season       string `json:"season"`
					SeasonYear   int    `json:"seasonYear"`
					Format       string `json:"format"`
					Source       string `json:"source"`
					Duration     int    `json:"duration"`
					AverageScore int    `json:"averageScore"`
					Popularity   int    `json:"popularity"`
					Favourites   int    `json:"favourites"`
					Rankings     []struct {
						Rank   int    `json:"rank"`
						Type   string `json:"type"`
						Season string `json:"season"`
						Year   int    `json:"year"`
					} `json:"rankings"`
					Studios struct {
						Nodes []struct {
							Name string `json:"name"`
						} `json:"nodes"`
					} `json:"studios"`
					Genres     []string `json:"genres"`
					Characters struct {
						Edges []struct {
							Node struct {
								Name struct {
									Full string `json:"full"`
								} `json:"name"`
								Image struct {
									Medium string `json:"medium"`
								} `json:"image"`
							} `json:"node"`
							Role        string `json:"role"`
							VoiceActors []struct {
								Name struct {
									Full string `json:"full"`
								} `json:"name"`
								Image struct {
									Medium string `json:"medium"`
								} `json:"image"`
							} `json:"voiceActors"`
						} `json:"edges"`
					} `json:"characters"`
					Staff struct {
						Edges []struct {
							Node struct {
								Name struct {
									Full string `json:"full"`
								} `json:"name"`
								Image struct {
									Medium string `json:"medium"`
								} `json:"image"`
							} `json:"node"`
							Role string `json:"role"`
						} `json:"edges"`
					} `json:"staff"`
					Relations struct {
						Edges []struct {
							Node struct {
								Title struct {
									Romaji  string `json:"romaji"`
									English string `json:"english"`
								} `json:"title"`
								CoverImage struct {
									Medium string `json:"medium"`
								} `json:"coverImage"`
							} `json:"node"`
							RelationType string `json:"relationType"`
						} `json:"edges"`
					} `json:"relations"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&anilistResp); err != nil {
		return nil, err
	}

	if len(anilistResp.Data.Page.Media) == 0 {
		return nil, fmt.Errorf("anime not found")
	}

	media := anilistResp.Data.Page.Media[0]

	// Convert to our format
	enhanced := &EnhancedAnimeData{
		AlternativeTitles: models.AlternativeTitles{
			Synonyms: media.Synonyms,
			Japanese: media.Title.Native,
			English:  media.Title.English,
		},
		Information: models.AnimeInformation{
			Episodes:  media.Episodes,
			Status:    media.Status,
			Aired:     formatAirDate(media.StartDate, media.EndDate),
			Premiered: fmt.Sprintf("%s %d", media.Season, media.SeasonYear),
			Studios:   extractStudioNames(media.Studios.Nodes),
			Source:    media.Source,
			Duration:  fmt.Sprintf("%d min", media.Duration),
		},
		Statistics: models.AnimeStatistics{
			Score:      float64(media.AverageScore) / 10.0,
			Ranked:     extractRanking(media.Rankings),
			Popularity: media.Popularity,
			Favorites:  media.Favourites,
		},
	}

	// Convert characters
	for _, edge := range media.Characters.Edges {
		char := models.Character{
			Name:     edge.Node.Name.Full,
			Role:     edge.Role,
			ImageUrl: edge.Node.Image.Medium,
		}
		if len(edge.VoiceActors) > 0 {
			char.VoiceActor = edge.VoiceActors[0].Name.Full
			char.VAImageUrl = edge.VoiceActors[0].Image.Medium
		}
		enhanced.Characters = append(enhanced.Characters, char)
	}

	// Convert staff
	for _, edge := range media.Staff.Edges {
		staff := models.StaffMember{
			Name:     edge.Node.Name.Full,
			Role:     edge.Role,
			ImageUrl: edge.Node.Image.Medium,
		}
		enhanced.Staff = append(enhanced.Staff, staff)
	}

	// Convert related anime
	for _, edge := range media.Relations.Edges {
		title := edge.Node.Title.English
		if title == "" {
			title = edge.Node.Title.Romaji
		}
		related := models.RelatedAnime{
			Name:         title,
			RelationType: edge.RelationType,
			ImageUrl:     edge.Node.CoverImage.Medium,
		}
		enhanced.Related = append(enhanced.Related, related)
	}

	// Get theme songs
	themes, err := GetAnimeThemes(animeName)
	if err == nil {
		enhanced.Themes = models.AnimeThemes{
			Openings: themes.Openings,
			Endings:  themes.Endings,
		}
	}

	return enhanced, nil
}

func formatAirDate(start, end struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}) string {
	if start.Year == 0 {
		return "Unknown"
	}

	startStr := fmt.Sprintf("%d-%02d-%02d", start.Year, start.Month, start.Day)
	if end.Year == 0 {
		return startStr + " to ?"
	}

	endStr := fmt.Sprintf("%d-%02d-%02d", end.Year, end.Month, end.Day)
	return startStr + " to " + endStr
}

func extractStudioNames(studios []struct {
	Name string `json:"name"`
}) []string {
	var names []string
	for _, studio := range studios {
		names = append(names, studio.Name)
	}
	return names
}

func extractRanking(rankings []struct {
	Rank   int    `json:"rank"`
	Type   string `json:"type"`
	Season string `json:"season"`
	Year   int    `json:"year"`
}) int {
	for _, ranking := range rankings {
		if ranking.Type == "RATED" {
			return ranking.Rank
		}
	}
	return 0
}