package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	model "github.com/Flack74/mongoapi/models"
	"github.com/Flack74/mongoapi/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// sendJSONResponse sends a standardized JSON response
func sendJSONResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: success,
		Message: message,
		Data:    data,
		Error:   errorMsg,
	}

	json.NewEncoder(w).Encode(response)
}

// renderAnimeCards renders anime cards for HTMX responses
func renderAnimeCards(w http.ResponseWriter, animes []primitive.M) {
	renderAnimeCardsWithLayout(w, animes, "grid")
}

// renderAnimeCardsWithLayout renders anime cards with specified layout
func renderAnimeCardsWithLayout(w http.ResponseWriter, animes []primitive.M, layout string) {
	if len(animes) == 0 {
		fmt.Fprintf(w, `<div class="col-span-full text-center py-12 text-gray-400">
			<p class="text-xl mb-2">😢 No anime found</p>
			<p>Try adjusting your search or filters</p>
		</div>`)
		return
	}
	
	for _, anime := range animes {
		name := ""
		imageUrl := "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjI2NyIgdmlld0JveD0iMCAwIDIwMCAyNjciIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjY3IiBmaWxsPSIjMzc0MTUxIi8+Cjx0ZXh0IHg9IjEwMCIgeT0iMTMzIiBmaWxsPSIjNkI3MjgwIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNCI+Tm8gSW1hZ2U8L3RleHQ+Cjwvc3ZnPg=="
		score := 0
		status := ""
		genres := []string{}
		year := 0
		season := ""
		animeType := ""
		
		// Extract data with type checking
		if n, ok := anime["name"].(string); ok {
			name = n
		}
		if img, ok := anime["imageUrl"].(string); ok && img != "" {
			imageUrl = img
		}
		// Handle different score types
		if s, ok := anime["score"].(int32); ok {
			score = int(s)
		} else if s, ok := anime["score"].(int64); ok {
			score = int(s)
		} else if s, ok := anime["score"].(int); ok {
			score = s
		} else if s, ok := anime["score"].(float64); ok {
			score = int(s)
		}
		if st, ok := anime["status"].(string); ok {
			status = st
		}
		if tp, ok := anime["type"].(string); ok {
			animeType = tp
		}
		// Handle year
		if y, ok := anime["year"].(int32); ok {
			year = int(y)
		} else if y, ok := anime["year"].(int64); ok {
			year = int(y)
		} else if y, ok := anime["year"].(int); ok {
			year = y
		} else if y, ok := anime["year"].(float64); ok {
			year = int(y)
		}
		if s, ok := anime["season"].(string); ok {
			season = s
		}
		// Handle genres
		if g, ok := anime["genre"].(primitive.A); ok {
			for _, genre := range g {
				if genreStr, ok := genre.(string); ok {
					genres = append(genres, genreStr)
				}
			}
		}
		
		genreStr := strings.Join(genres, ", ")
		if len(genreStr) > 25 {
			genreStr = genreStr[:25] + "..."
		}
		
		// Build year/season string
		yearSeasonStr := ""
		if year > 0 {
			yearSeasonStr = fmt.Sprintf("%d", year)
			if season != "" {
				yearSeasonStr += " " + season
			}
		} else if season != "" {
			yearSeasonStr = season
		}
		
		// Different layouts for grid vs horizontal
		className := "bg-gray-800 rounded-lg overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105 cursor-pointer"
		if layout == "horizontal" {
			className = "flex-none w-48 " + className
		}
		
		fmt.Fprintf(w, `
		<div class="%s"
		     onclick="fetch('/api/anime/%s', {headers: {'HX-Request': 'true'}}).then(r => r.text()).then(html => {document.getElementById('modal-content').innerHTML = html; showModal();})">
		    <div class="aspect-[3/4] bg-gray-700 relative overflow-hidden">
		        <img src="%s" alt="%s" class="w-full h-full object-cover" 
		             onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjI2NyIgdmlld0JveD0iMCAwIDIwMCAyNjciIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjY3IiBmaWxsPSIjMzc0MTUxIi8+Cjx0ZXh0IHg9IjEwMCIgeT0iMTMzIiBmaWxsPSIjNkI3MjgwIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNCI+Tm8gSW1hZ2U8L3RleHQ+Cjwvc3ZnPg=='">
		        <div class="absolute top-2 right-2 bg-anime-blue text-white px-2 py-1 rounded text-sm font-bold">
		            %d/10
		        </div>
		        <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-3">
		            <h3 class="text-white font-semibold text-sm mb-1 truncate" title="%s">%s</h3>
		            <p class="text-gray-300 text-xs mb-1">%s • %s</p>
		            <p class="text-gray-400 text-xs truncate">%s</p>
		            <p class="text-gray-500 text-xs mt-1">%s</p>
		        </div>
		    </div>
		</div>`, 
			className,
			strings.ReplaceAll(strings.ToLower(name), " ", "-"), 
			imageUrl, name, score, name, name, animeType, status, genreStr, yearSeasonStr)
	}
}

// renderAnimeModal renders detailed anime information for modal
func renderAnimeModal(w http.ResponseWriter, anime *model.Anime) {
	genreStr := strings.Join(anime.Genre, ", ")
	if genreStr == "" {
		genreStr = "Not specified"
	}
	
	progressText := "Not specified"
	if anime.Progress.Total > 0 {
		if anime.Progress.Watched > 0 {
			progressText = fmt.Sprintf("%d/%d episodes", anime.Progress.Watched, anime.Progress.Total)
		} else {
			progressText = fmt.Sprintf("%d episodes", anime.Progress.Total)
		}
	}
	
	yearSeason := ""
	if anime.Year > 0 {
		yearSeason = fmt.Sprintf("%d", anime.Year)
		if anime.Season != "" {
			yearSeason += fmt.Sprintf(" %s", anime.Season)
		}
	}
	
	fmt.Fprintf(w, `
	<div class="flex flex-col md:flex-row gap-6">
		<div class="flex-shrink-0">
			<img src="%s" alt="%s" class="w-48 h-64 object-cover rounded-lg" 
			     onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjI2NyIgdmlld0JveD0iMCAwIDIwMCAyNjciIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjY3IiBmaWxsPSIjMzc0MTUxIi8+Cjx0ZXh0IHg9IjEwMCIgeT0iMTMzIiBmaWxsPSIjNkI3MjgwIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNCI+Tm8gSW1hZ2U8L3RleHQ+Cjwvc3ZnPg=='">
		</div>
		<div class="flex-1">
			<h2 class="text-2xl font-bold text-white mb-4">%s</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
				<div>
					<span class="text-gray-400">Type:</span>
					<span class="text-white ml-2">%s</span>
				</div>
				<div>
					<span class="text-gray-400">Score:</span>
					<span class="text-anime-blue ml-2 font-bold">%d/10</span>
				</div>
				<div>
					<span class="text-gray-400">Status:</span>
					<span class="text-white ml-2 capitalize">%s</span>
				</div>
				<div>
					<span class="text-gray-400">Progress:</span>
					<span class="text-white ml-2">%s</span>
				</div>
				<div class="md:col-span-2">
					<span class="text-gray-400">Genres:</span>
					<span class="text-white ml-2">%s</span>
				</div>
				<div class="md:col-span-2">
					<span class="text-gray-400">Year/Season:</span>
					<span class="text-white ml-2">%s</span>
				</div>
			</div>
			<div class="mt-6">
				<h3 class="text-lg font-semibold text-white mb-2">Synopsis</h3>
				<p class="text-gray-300 leading-relaxed">%s</p>
			</div>
			<div class="mt-6 flex gap-3">
				<button class="bg-anime-blue hover:bg-blue-600 px-4 py-2 rounded-lg transition-colors"
				        hx-post="/api/admin/anime/%s/episode/increment" 
				        hx-target="#modal-content" 
				        hx-swap="outerHTML">
					+ Episode
				</button>
				<button class="bg-gray-600 hover:bg-gray-500 px-4 py-2 rounded-lg transition-colors"
				        hx-post="/api/admin/anime/%s/episode/decrement" 
				        hx-target="#modal-content" 
				        hx-swap="outerHTML">
					- Episode
				</button>
				<button class="bg-anime-purple hover:bg-purple-600 px-4 py-2 rounded-lg transition-colors"
				        hx-post="/api/admin/anime/%s/status/toggle" 
				        hx-target="#modal-content" 
				        hx-swap="outerHTML">
					Toggle Status
				</button>
			</div>
		</div>
	</div>`,
		anime.ImageUrl, anime.Name, anime.Name, anime.Type, anime.Score, 
		anime.Status, progressText, genreStr, yearSeason, anime.Notes,
		anime.ID.Hex(), anime.ID.Hex(), anime.ID.Hex())
}

func GetMyAllAnimesHandler(w http.ResponseWriter, r *http.Request) {
	allAnimes := services.GetAllAnimes()
	if allAnimes == nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to fetch animes")
		return
	}
	
	// Check if request wants HTML (HTMX)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderAnimeCards(w, allAnimes)
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Animes retrieved successfully", allAnimes, "")
}

func GetAnimeByNameHandler(w http.ResponseWriter, r *http.Request) {
	animeName := chi.URLParam(r, "animeName")
	if animeName == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Anime name is required")
		return
	}

	// Normalize anime name
	animeName = strings.ReplaceAll(animeName, "-", " ")
	animeName = strings.ReplaceAll(animeName, "_", " ")
	animeName = strings.ToLower(animeName)

	existingAnime, err := services.SearchAnimeByName(animeName)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Database error occurred")
		return
	}
	if existingAnime == nil {
		sendJSONResponse(w, http.StatusNotFound, false, "", nil, "Anime not found")
		return
	}

	// Check if request wants HTML (HTMX modal)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderAnimeModal(w, existingAnime)
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "Anime retrieved successfully", existingAnime, "")
}

func CreateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var anime model.Anime
	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	// Validate required fields
	if anime.Name == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Anime name is required")
		return
	}

	// Check if anime already exists
	existingAnime, err := services.FindAnimeByName(anime.Name)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Database error occurred")
		return
	}
	if existingAnime != nil {
		sendJSONResponse(w, http.StatusConflict, false, "", nil, "Anime with this name already exists")
		return
	}

	// Insert anime
	err = services.InsertOneAnime(anime)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to create anime")
		return
	}

	sendJSONResponse(w, http.StatusCreated, true, "Anime created successfully", anime, "")
}

func CreateMultipleAnimesHandler(w http.ResponseWriter, r *http.Request) {
	var animes []model.Anime
	if err := json.NewDecoder(r.Body).Decode(&animes); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	if len(animes) == 0 {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "No animes provided")
		return
	}

	insertedIDs, duplicates, err := services.InsertMultipleAnimes(animes)
	if err != nil {
		if len(duplicates) > 0 {
			sendJSONResponse(w, http.StatusPartialContent, false, "", map[string]interface{}{
				"duplicates": duplicates,
			}, "Some animes already exist or failed to insert")
		} else {
			sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to insert animes")
		}
		return
	}

	responseData := map[string]interface{}{
		"inserted_count": len(insertedIDs),
		"inserted_ids":   insertedIDs,
	}

	if len(duplicates) > 0 {
		responseData["duplicates"] = duplicates
		sendJSONResponse(w, http.StatusPartialContent, true, "Animes inserted with some duplicates skipped", responseData, "")
	} else {
		sendJSONResponse(w, http.StatusCreated, true, "All animes created successfully", responseData, "")
	}
}

func UpdateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	services.UpdateAnime(w, r)
}

func DeleteAnAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Anime ID is required")
		return
	}

	deleted := services.DeleteOneAnime(id)
	if !deleted {
		sendJSONResponse(w, http.StatusNotFound, false, "", nil, "Anime not found")
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "Anime deleted successfully", map[string]string{"deleted_id": id}, "")
}

func DeleteEveryAnimesHandler(w http.ResponseWriter, r *http.Request) {
	count := services.DeleteAllAnime()
	if count == 0 {
		sendJSONResponse(w, http.StatusNotFound, false, "", nil, "No animes found to delete")
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "All animes deleted successfully", map[string]int64{"deleted_count": count}, "")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, http.StatusOK, true, "API is healthy", map[string]string{
		"status":    "healthy",
		"version":   "3.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}, "")
}

func FilterAnimesHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	genre := r.URL.Query().Get("genre")
	year := r.URL.Query().Get("year")
	season := r.URL.Query().Get("season")
	format := r.URL.Query().Get("format")
	status := r.URL.Query().Get("status")

	filteredAnimes := services.FilterAnimes(search, genre, year, season, format, status)
	if filteredAnimes == nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to filter animes")
		return
	}
	
	// Check if request wants HTML (HTMX)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderAnimeCards(w, filteredAnimes)
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Animes filtered successfully", filteredAnimes, "")
}

func GetTrendingAnimesHandler(w http.ResponseWriter, r *http.Request) {
	trendingAnimes := services.GetTrendingAnimes()
	if trendingAnimes == nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to fetch trending animes")
		return
	}
	
	// Check if request wants HTML (HTMX)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderAnimeCardsWithLayout(w, trendingAnimes, "horizontal")
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Trending animes retrieved successfully", trendingAnimes, "")
}

func GetPopularAnimesHandler(w http.ResponseWriter, r *http.Request) {
	popularAnimes := services.GetPopularAnimes()
	if popularAnimes == nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to fetch popular animes")
		return
	}
	
	// Check if request wants HTML (HTMX)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderAnimeCardsWithLayout(w, popularAnimes, "horizontal")
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Popular animes retrieved successfully", popularAnimes, "")
}

func ServeFrontendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en" class="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AnimeVerse - Modern Anime Database</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            darkMode: 'class',
            theme: {
                extend: {
                    colors: {
                        'anime-blue': '#3B82F6',
                        'anime-purple': '#8B5CF6',
                    }
                }
            }
        }
    </script>
    <style>
        .scrollbar-hide {
            -ms-overflow-style: none;
            scrollbar-width: none;
        }
        .scrollbar-hide::-webkit-scrollbar {
            display: none;
        }
        .htmx-indicator {
            display: none;
        }
        .htmx-request .htmx-indicator {
            display: block;
        }
    </style>
</head>
<body class="bg-gray-900 text-white min-h-screen">
    <!-- Header -->
    <header class="bg-gray-800 shadow-lg">
        <div class="container mx-auto px-4 py-6">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-4">
                    <div class="text-3xl font-bold bg-gradient-to-r from-anime-blue to-anime-purple bg-clip-text text-transparent">
                        🌸 AnimeVerse
                    </div>
                    <span class="text-sm text-gray-400">v3.0</span>
                </div>
                <nav class="hidden md:flex space-x-4">
                    <a href="#" class="text-gray-300 hover:text-white transition-colors">Home</a>
                    <a href="/api/animes" class="text-gray-300 hover:text-white transition-colors">API</a>
                    <button class="bg-anime-blue hover:bg-blue-600 px-3 py-1 rounded text-sm transition-colors"
                            hx-post="/api/admin/import/trending" 
                            hx-trigger="click"
                            hx-confirm="Import 25 trending anime from MyAnimeList?">
                        Import Trending
                    </button>
                    <button class="bg-anime-purple hover:bg-purple-600 px-3 py-1 rounded text-sm transition-colors"
                            hx-post="/api/admin/import/seasonal?year=2024&season=winter" 
                            hx-trigger="click"
                            hx-confirm="Import 2024 Winter seasonal anime?">
                        Import Seasonal
                    </button>
                </nav>
            </div>
        </div>
    </header>

    <!-- Filters -->
    <section class="bg-gray-800 border-b border-gray-700">
        <div class="container mx-auto px-4 py-6">
            <form id="filter-form" class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4">
                <input type="text" name="search" id="search-input" placeholder="Search anime..." 
                       class="bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-anime-blue">
                
                <select name="genre" id="genre-select" class="bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-anime-blue">
                    <option value="">All Genres</option>
                    <option value="Action">Action</option>
                    <option value="Adventure">Adventure</option>
                    <option value="Comedy">Comedy</option>
                    <option value="Drama">Drama</option>
                    <option value="Fantasy">Fantasy</option>
                    <option value="Romance">Romance</option>
                    <option value="Thriller">Thriller</option>
                    <option value="Mystery">Mystery</option>
                    <option value="Supernatural">Supernatural</option>
                    <option value="Horror">Horror</option>
                    <option value="School">School</option>
                </select>
                
                <select name="year" id="year-select" class="bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-anime-blue">
                    <option value="">All Years</option>
                    <option value="2024">2024</option>
                    <option value="2023">2023</option>
                    <option value="2022">2022</option>
                    <option value="2021">2021</option>
                    <option value="2020">2020</option>
                    <option value="2019">2019</option>
                    <option value="2018">2018</option>
                    <option value="2017">2017</option>
                    <option value="2016">2016</option>
                </select>
                
                <select name="season" id="season-select" class="bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-anime-blue">
                    <option value="">All Seasons</option>
                    <option value="Winter">Winter</option>
                    <option value="Spring">Spring</option>
                    <option value="Summer">Summer</option>
                    <option value="Fall">Fall</option>
                </select>
                
                <select name="format" id="format-select" class="bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-anime-blue">
                    <option value="">All Formats</option>
                    <option value="TV">TV</option>
                    <option value="Movie">Movie</option>
                    <option value="ONA">ONA</option>
                    <option value="OVA">OVA</option>
                </select>
                
                <select name="status" id="status-select" class="bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-anime-blue">
                    <option value="">All Status</option>
                    <option value="watching">Watching</option>
                    <option value="completed">Completed</option>
                    <option value="on-hold">On Hold</option>
                    <option value="dropped">Dropped</option>
                    <option value="plan-to-watch">Plan to Watch</option>
                </select>
            </form>
            
            <!-- Loading Indicator -->
            <div id="loading" class="hidden mt-4 text-center">
                <div class="inline-flex items-center px-4 py-2 bg-anime-blue rounded-lg">
                    <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Searching...
                </div>
            </div>
        </div>
    </section>

    <!-- Main Content -->
    <main class="container mx-auto px-4 py-8">
        <!-- Trending Section -->
        <section class="mb-12">
            <h2 class="text-2xl font-bold mb-6 text-anime-blue">🔥 Trending Now</h2>
            <div class="relative">
                <div id="trending-grid" class="flex gap-6 overflow-x-auto pb-4 scrollbar-hide"
                     hx-get="/api/animes/trending" hx-trigger="load" hx-swap="innerHTML">
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                </div>
            </div>
        </section>

        <!-- Popular Section -->
        <section class="mb-12">
            <h2 class="text-2xl font-bold mb-6 text-anime-purple">⭐ Popular This Season</h2>
            <div class="relative">
                <div id="popular-grid" class="flex gap-6 overflow-x-auto pb-4 scrollbar-hide"
                     hx-get="/api/animes/popular" hx-trigger="load" hx-swap="innerHTML">
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                    <div class="flex-none w-48 animate-pulse bg-gray-700 rounded-lg h-72"></div>
                </div>
            </div>
        </section>

        <!-- All Anime Grid -->
        <section>
            <h2 class="text-2xl font-bold mb-6">📚 All Anime</h2>
            <div id="anime-grid" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6"
                 hx-get="/api/animes" hx-trigger="load" hx-swap="innerHTML">
                <div class="animate-pulse bg-gray-700 rounded-lg h-80"></div>
                <div class="animate-pulse bg-gray-700 rounded-lg h-80"></div>
                <div class="animate-pulse bg-gray-700 rounded-lg h-80"></div>
            </div>
        </section>
    </main>

    <!-- Modal -->
    <div id="modal" class="fixed inset-0 bg-black bg-opacity-50 hidden items-center justify-center z-50">
        <div class="bg-gray-800 rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
            <div id="modal-content"></div>
            <button onclick="document.getElementById('modal').classList.add('hidden')" 
                    class="mt-4 bg-anime-blue hover:bg-blue-600 px-4 py-2 rounded-lg transition-colors">
                Close
            </button>
        </div>
    </div>

    <script>
        // Show modal function
        function showModal() {
            document.getElementById('modal').classList.remove('hidden');
            document.getElementById('modal').classList.add('flex');
        }
        
        // Filter function
        function applyFilters() {
            const form = document.getElementById('filter-form');
            const formData = new FormData(form);
            const params = new URLSearchParams();
            
            // Add non-empty values to params
            for (let [key, value] of formData.entries()) {
                if (value.trim() !== '') {
                    params.append(key, value);
                }
            }
            
            // Show loading
            document.getElementById('loading').classList.remove('hidden');
            document.getElementById('anime-grid').style.opacity = '0.5';
            
            // Make request
            fetch('/api/animes/filter?' + params.toString(), {
                headers: {
                    'HX-Request': 'true'
                }
            })
            .then(response => response.text())
            .then(html => {
                document.getElementById('anime-grid').innerHTML = html;
                document.getElementById('anime-grid').style.opacity = '1';
                document.getElementById('loading').classList.add('hidden');
            })
            .catch(error => {
                console.error('Filter error:', error);
                document.getElementById('anime-grid').style.opacity = '1';
                document.getElementById('loading').classList.add('hidden');
            });
        }
        
        // Debounced search
        let searchTimeout;
        function debounceSearch() {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(applyFilters, 500);
        }
        
        // Event listeners
        document.addEventListener('DOMContentLoaded', function() {
            // Search input
            document.getElementById('search-input').addEventListener('input', debounceSearch);
            
            // Filter selects
            document.getElementById('genre-select').addEventListener('change', applyFilters);
            document.getElementById('year-select').addEventListener('change', applyFilters);
            document.getElementById('season-select').addEventListener('change', applyFilters);
            document.getElementById('format-select').addEventListener('change', applyFilters);
            document.getElementById('status-select').addEventListener('change', applyFilters);
            
            // Load initial data
            applyFilters();
        });
    </script>
</body>
</html>
	`))
}

func IncrementEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	anime, err := services.IncrementEpisode(id)
	if err != nil {
		http.Error(w, "Failed to update episode", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	renderAnimeModal(w, anime)
}

func DecrementEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	anime, err := services.DecrementEpisode(id)
	if err != nil {
		http.Error(w, "Failed to update episode", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	renderAnimeModal(w, anime)
}

func ToggleStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	anime, err := services.ToggleStatus(id)
	if err != nil {
		http.Error(w, "Failed to toggle status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	renderAnimeModal(w, anime)
}

func ImportTrendingHandler(w http.ResponseWriter, r *http.Request) {
	count, err := services.ImportTrendingAnime()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to import anime: "+err.Error())
		return
	}
	sendJSONResponse(w, http.StatusOK, true, fmt.Sprintf("Imported %d trending anime", count), nil, "")
}

func ImportSeasonalHandler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	season := r.URL.Query().Get("season")
	
	if year == "" || season == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Year and season parameters required")
		return
	}
	
	count, err := services.ImportSeasonalAnime(year, season)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to import seasonal anime: "+err.Error())
		return
	}
	sendJSONResponse(w, http.StatusOK, true, fmt.Sprintf("Imported %d seasonal anime", count), nil, "")
}

func ServeHomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>AnimeVerse API v3.0</title>
			<style>
				body {
					font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
					background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
					color: #fff;
					text-align: center;
					padding: 50px;
					min-height: 100vh;
					margin: 0;
				}
				h1 {
					font-size: 3rem;
					color: #fff;
					margin-bottom: 10px;
					text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
				}
				.version {
					font-size: 1rem;
					color: #ffeb3b;
					margin-bottom: 30px;
				}
				.container {
					max-width: 700px;
					margin: auto;
					padding: 40px;
					background: rgba(255, 255, 255, 0.1);
					border-radius: 20px;
					box-shadow: 0 8px 32px rgba(0,0,0,0.3);
					backdrop-filter: blur(10px);
					border: 1px solid rgba(255,255,255,0.2);
				}
				p {
					font-size: 1.3rem;
					margin: 20px 0;
					color: #f0f0f0;
					line-height: 1.6;
				}
				.features {
					display: grid;
					grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
					gap: 20px;
					margin-top: 30px;
				}
				.feature {
					padding: 20px;
					background: rgba(255,255,255,0.1);
					border-radius: 10px;
					border: 1px solid rgba(255,255,255,0.2);
				}
				.api-link {
					display: inline-block;
					margin-top: 20px;
					padding: 12px 24px;
					background: #4CAF50;
					color: white;
					text-decoration: none;
					border-radius: 25px;
					transition: background 0.3s;
				}
				.api-link:hover {
					background: #45a049;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>🌸 AnimeVerse API</h1>
				<div class="version">v3.0 - Production Ready</div>
				<p>Your ultimate RESTful API for managing and exploring anime collections!</p>
				<div class="features">
					<div class="feature">
						<h3>⚡ Fast & Lightweight</h3>
						<p>Built with Chi router for optimal performance</p>
					</div>
					<div class="feature">
						<h3>🔒 CORS Enabled</h3>
						<p>Ready for web applications</p>
					</div>
					<div class="feature">
						<h3>📊 Standardized</h3>
						<p>Consistent JSON responses</p>
					</div>
				</div>
				<a href="/api/animes" class="api-link">View All Anime →</a>
			</div>
		</body>
		</html>
	`))
}
