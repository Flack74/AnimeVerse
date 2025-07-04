package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

// GetFastBrowseHandler handles fast browse requests
func GetFastBrowseHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	animes, err := fetchEnhancedAnimeData("", page)
	if err != nil {
		http.Error(w, "Failed to load anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    animes,
		"page":    page,
	})
}

// GetFastTopRatedHandler handles fast top-rated requests
func GetFastTopRatedHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	animes, err := fetchTopRatedAnimeData(page)
	if err != nil {
		http.Error(w, "Failed to load top rated anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    animes,
		"page":    page,
	})
}

// GetFastSearchHandler handles fast search requests
func GetFastSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query required", http.StatusBadRequest)
		return
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	animes, err := fetchEnhancedAnimeData(query, page)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    animes,
		"query":   query,
		"page":    page,
	})
}

// fetchEnhancedAnimeData fetches anime with high-quality images and metadata
func fetchEnhancedAnimeData(search string, page int) ([]map[string]interface{}, error) {
	graphqlQuery := `
	query ($page: Int, $search: String) {
		Page(page: $page, perPage: 24) {
			media(search: $search, type: ANIME, sort: POPULARITY_DESC) {
				id
				title {
					romaji
					english
					native
				}
				description(asHtml: false)
				genres
				episodes
				status
				format
				startDate {
					year
				}
				averageScore
				coverImage {
					extraLarge
					large
					medium
				}
				bannerImage
				studios {
					nodes {
						name
					}
				}
				trailer {
					id
					site
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"page": page,
	}
	if search != "" {
		variables["search"] = search
	}

	requestBody := map[string]interface{}{
		"query":     graphqlQuery,
		"variables": variables,
	}

	jsonData, _ := json.Marshal(requestBody)
	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var animes []map[string]interface{}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if page, ok := data["Page"].(map[string]interface{}); ok {
			if media, ok := page["media"].([]interface{}); ok {
				for _, item := range media {
					if anime, ok := item.(map[string]interface{}); ok {
						// Process anime data
						processedAnime := processAnimeData(anime)
						animes = append(animes, processedAnime)
					}
				}
			}
		}
	}

	return animes, nil
}

// fetchTopRatedAnimeData fetches top-rated anime
func fetchTopRatedAnimeData(page int) ([]map[string]interface{}, error) {
	graphqlQuery := `
	query ($page: Int) {
		Page(page: $page, perPage: 24) {
			media(type: ANIME, sort: SCORE_DESC) {
				id
				title {
					romaji
					english
					native
				}
				description(asHtml: false)
				genres
				episodes
				status
				format
				startDate {
					year
				}
				averageScore
				coverImage {
					extraLarge
					large
					medium
				}
				bannerImage
				studios {
					nodes {
						name
					}
				}
			}
		}
	}`

	requestBody := map[string]interface{}{
		"query": graphqlQuery,
		"variables": map[string]interface{}{
			"page": page,
		},
	}

	jsonData, _ := json.Marshal(requestBody)
	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var animes []map[string]interface{}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if page, ok := data["Page"].(map[string]interface{}); ok {
			if media, ok := page["media"].([]interface{}); ok {
				for _, item := range media {
					if anime, ok := item.(map[string]interface{}); ok {
						processedAnime := processAnimeData(anime)
						animes = append(animes, processedAnime)
					}
				}
			}
		}
	}

	return animes, nil
}

// processAnimeData processes raw anime data into clean format
func processAnimeData(anime map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})

	// Title
	if title, ok := anime["title"].(map[string]interface{}); ok {
		if english, ok := title["english"].(string); ok && english != "" {
			processed["name"] = english
		} else if romaji, ok := title["romaji"].(string); ok {
			processed["name"] = romaji
		}
	}

	// Basic info
	processed["synopsis"] = anime["description"]
	processed["genre"] = anime["genres"]
	processed["year"] = extractYear(anime["startDate"])
	processed["type"] = anime["format"]

	// Score
	if score, ok := anime["averageScore"].(float64); ok {
		processed["score"] = score / 10.0
	}

	// Images - prioritize high quality
	if coverImage, ok := anime["coverImage"].(map[string]interface{}); ok {
		if extraLarge, ok := coverImage["extraLarge"].(string); ok && extraLarge != "" {
			processed["imageUrl"] = extraLarge
		} else if large, ok := coverImage["large"].(string); ok && large != "" {
			processed["imageUrl"] = large
		}
	}

	if bannerImage, ok := anime["bannerImage"].(string); ok {
		processed["bannerUrl"] = bannerImage
	}

	// Information
	info := map[string]interface{}{
		"episodes": anime["episodes"],
		"status":   anime["status"],
	}

	// Studios
	if studios, ok := anime["studios"].(map[string]interface{}); ok {
		if nodes, ok := studios["nodes"].([]interface{}); ok {
			var studioNames []string
			for _, node := range nodes {
				if studio, ok := node.(map[string]interface{}); ok {
					if name, ok := studio["name"].(string); ok {
						studioNames = append(studioNames, name)
					}
				}
			}
			info["studios"] = studioNames
		}
	}

	processed["information"] = info

	return processed
}

// extractYear extracts year from startDate
func extractYear(startDate interface{}) interface{} {
	if date, ok := startDate.(map[string]interface{}); ok {
		return date["year"]
	}
	return nil
}