package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	model "github.com/Flack74/mongoapi/models"
)

type JikanAnime struct {
	Title    string `json:"title"`
	Type     string `json:"type"`
	Score    float64 `json:"score"`
	Year     int    `json:"year"`
	Season   string `json:"season"`
	Genres   []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Images struct {
		JPG struct {
			ImageURL string `json:"image_url"`
		} `json:"jpg"`
	} `json:"images"`
	Synopsis string `json:"synopsis"`
}

type JikanResponse struct {
	Data []JikanAnime `json:"data"`
}

func ImportTrendingAnime() (int, error) {
	resp, err := http.Get("https://api.jikan.moe/v4/top/anime?limit=25")
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

func ImportSeasonalAnime(year, season string) (int, error) {
	url := fmt.Sprintf("https://api.jikan.moe/v4/seasons/%s/%s?limit=25", year, strings.ToLower(season))
	resp, err := http.Get(url)
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

func importAnimeList(jikanAnimes []JikanAnime) (int, error) {
	var animes []model.Anime
	imported := 0

	for _, ja := range jikanAnimes {
		// Check if already exists
		existing, _ := FindAnimeByName(ja.Title)
		if existing != nil {
			continue
		}

		// Convert genres
		var genres []string
		for _, g := range ja.Genres {
			genres = append(genres, g.Name)
		}

		// Convert season
		season := model.Season("")
		switch strings.ToLower(ja.Season) {
		case "winter":
			season = model.Winter
		case "spring":
			season = model.Spring
		case "summer":
			season = model.Summer
		case "fall":
			season = model.Fall
		}

		anime := model.Anime{
			Name:      ja.Title,
			Type:      model.AnimeType(ja.Type),
			Score:     int(ja.Score),
			Status:    "plan-to-watch",
			Genre:     genres,
			Notes:     truncateString(ja.Synopsis, 500),
			Year:      ja.Year,
			Season:    season,
			ImageUrl:  ja.Images.JPG.ImageURL,
			Progress: model.Progress{
				Watched: 0,
				Total:   0,
			},
		}

		animes = append(animes, anime)
		imported++

		// Rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	if len(animes) > 0 {
		_, _, err := InsertMultipleAnimes(animes)
		if err != nil {
			return 0, err
		}
	}

	return imported, nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}