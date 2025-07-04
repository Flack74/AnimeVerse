package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// SimpleBrowseHandler - Direct AniList proxy for browse page
func SimpleBrowseHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	genre := r.URL.Query().Get("genre")
	year := r.URL.Query().Get("year")
	search := r.URL.Query().Get("search")

	// Build AniList query
	query := fmt.Sprintf(`{
		Page(page: %d, perPage: 50) {
			media(type: ANIME, sort: POPULARITY_DESC%s) {
				id
				title { romaji english }
				coverImage { extraLarge large }
				averageScore
				startDate { year }
				genres
				format
				status
				episodes
			}
		}
	}`, page, buildFilters(genre, year, search))

	// Make request to AniList
	requestBody := map[string]interface{}{
		"query": query,
	}

	jsonData, _ := json.Marshal(requestBody)
	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to fetch anime data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, "Failed to parse anime data", http.StatusInternalServerError)
		return
	}

	// Extract media data
	var animes []interface{}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if page, ok := data["Page"].(map[string]interface{}); ok {
			if media, ok := page["media"].([]interface{}); ok {
				animes = media
			}
		}
	}

	// Return the data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    animes,
		"source":  "anilist",
	})
}

func buildFilters(genre, year, search string) string {
	filters := ""
	if search != "" {
		filters += fmt.Sprintf(`, search: "%s"`, search)
	}
	if genre != "" {
		filters += fmt.Sprintf(`, genre_in: ["%s"]`, genre)
	}
	if year != "" {
		filters += fmt.Sprintf(`, seasonYear: %s`, year)
	}
	return filters
}