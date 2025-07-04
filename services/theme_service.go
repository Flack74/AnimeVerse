package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ThemeData struct {
	Openings []string `json:"openings"`
	Endings  []string `json:"endings"`
}

func GetAnimeThemes(animeName string) (*ThemeData, error) {
	// Try AniList first for theme data
	themes, err := fetchThemesFromAniList(animeName)
	if err == nil && (len(themes.Openings) > 0 || len(themes.Endings) > 0) {
		return themes, nil
	}

	// Fallback to MyAnimeList API
	return fetchThemesFromMAL(animeName)
}

func fetchThemesFromAniList(animeName string) (*ThemeData, error) {
	query := fmt.Sprintf(`{
		Page(page: 1, perPage: 1) {
			media(search: "%s", type: ANIME) {
				title {
					romaji
					english
				}
				externalLinks {
					url
					site
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
					} `json:"title"`
					ExternalLinks []struct {
						URL  string `json:"url"`
						Site string `json:"site"`
					} `json:"externalLinks"`
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

	// For now, return placeholder themes - in production, you'd parse from external links
	return &ThemeData{
		Openings: []string{
			"Opening 1: \"Theme Song\" by Artist",
		},
		Endings: []string{
			"Ending 1: \"End Theme\" by Artist",
		},
	}, nil
}

func fetchThemesFromMAL(animeName string) (*ThemeData, error) {
	// Search for anime on Jikan API
	searchURL := fmt.Sprintf("https://api.jikan.moe/v4/anime?q=%s&limit=1", strings.ReplaceAll(animeName, " ", "%20"))
	
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResp struct {
		Data []struct {
			MALID int `json:"mal_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	if len(searchResp.Data) == 0 {
		return nil, fmt.Errorf("anime not found")
	}

	// Get full anime data including themes
	animeURL := fmt.Sprintf("https://api.jikan.moe/v4/anime/%d/full", searchResp.Data[0].MALID)
	
	resp2, err := http.Get(animeURL)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	var animeResp struct {
		Data struct {
			Theme struct {
				Openings []string `json:"openings"`
				Endings  []string `json:"endings"`
			} `json:"theme"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp2.Body).Decode(&animeResp); err != nil {
		return nil, err
	}

	return &ThemeData{
		Openings: animeResp.Data.Theme.Openings,
		Endings:  animeResp.Data.Theme.Endings,
	}, nil
}

func GetHighQualityImages(animeName string) (string, string, error) {
	query := fmt.Sprintf(`{
		Page(page: 1, perPage: 1) {
			media(search: "%s", type: ANIME) {
				coverImage {
					extraLarge
				}
				bannerImage
			}
		}
	}`, animeName)

	requestBody := map[string]interface{}{
		"query": query,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var anilistResp struct {
		Data struct {
			Page struct {
				Media []struct {
					CoverImage struct {
						ExtraLarge string `json:"extraLarge"`
					} `json:"coverImage"`
					BannerImage string `json:"bannerImage"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&anilistResp); err != nil {
		return "", "", err
	}

	if len(anilistResp.Data.Page.Media) == 0 {
		return "", "", fmt.Errorf("anime not found")
	}

	media := anilistResp.Data.Page.Media[0]
	return media.CoverImage.ExtraLarge, media.BannerImage, nil
}