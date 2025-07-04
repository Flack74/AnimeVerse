package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"animeverse/middleware"
	model "animeverse/models"
	"animeverse/services"
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
		<div class="bg-white rounded-2xl shadow-lg overflow-hidden anime-card cursor-pointer"
		     onclick="showAnimeModal('%s')">
		    <div class="aspect-[3/4] bg-gray-200 relative overflow-hidden">
		        <img src="%s" alt="%s" class="w-full h-full object-cover" 
		             onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjI2NyIgdmlld0JveD0iMCAwIDIwMCAyNjciIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjY3IiBmaWxsPSIjRjNGNEY2Ii8+Cjx0ZXh0IHg9IjEwMCIgeT0iMTMzIiBmaWxsPSIjOUI5QkE0IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNCI+Tm8gSW1hZ2U8L3RleHQ+Cjwvc3ZnPg=='">
		        <div class="absolute top-3 right-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white px-2 py-1 rounded-full text-xs font-bold">
		            %d★
		        </div>
		    </div>
		    <div class="p-4">
		        <h3 class="font-bold text-gray-800 text-sm mb-2 truncate" title="%s">%s</h3>
		        <div class="flex items-center justify-between text-xs text-gray-500">
		            <span>%s</span>
		            <span class="bg-gray-100 px-2 py-1 rounded-full">%s</span>
		        </div>
		        <p class="text-gray-400 text-xs mt-2 truncate">%s</p>
		    </div>
		</div>`, 
			strings.ReplaceAll(strings.ToLower(name), " ", "-"), 
			imageUrl, name, score, name, name, animeType, status, genreStr)
	}
}

// renderAnimeModal renders detailed anime information for modal with tabs
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
	<div class="flex justify-between items-center mb-8">
		<h1 class="text-4xl font-bold text-gray-800">%s</h1>
		<button onclick="closeModal()" class="text-gray-400 hover:text-gray-600 text-3xl">&times;</button>
	</div>
	
	<!-- Hero Section -->
	<div class="relative h-64 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-2xl mb-8 overflow-hidden">
		<img src="%s" alt="%s banner" class="w-full h-full object-cover opacity-30" 
		     onerror="this.style.display='none'">
		<div class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent"></div>
		<div class="absolute bottom-6 left-6 flex items-end space-x-6">
			<img src="%s" alt="%s" class="w-32 h-44 object-cover rounded-2xl shadow-2xl" 
			     onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjI2NyIgdmlld0JveD0iMCAwIDIwMCAyNjciIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjY3IiBmaWxsPSIjRjNGNEY2Ii8+Cjx0ZXh0IHg9IjEwMCIgeT0iMTMzIiBmaWxsPSIjOUI5QkE0IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNCI+Tm8gSW1hZ2U8L3RleHQ+Cjwvc3ZnPg=='">
			<div class="text-white">
				<div class="flex items-center space-x-3 mb-3">
					<span class="bg-yellow-400 text-black px-3 py-1 rounded-full text-sm font-bold">%d★</span>
					<span class="bg-white/20 backdrop-blur px-3 py-1 rounded-full text-sm font-medium">16+</span>
					<span class="bg-white/20 backdrop-blur px-3 py-1 rounded-full text-sm font-medium">%s</span>
					<span class="bg-white/20 backdrop-blur px-3 py-1 rounded-full text-sm font-medium">23 min</span>
				</div>
				<div class="flex flex-wrap gap-2">
					<span class="bg-indigo-500 px-2 py-1 rounded text-xs">Action</span>
					<span class="bg-purple-500 px-2 py-1 rounded text-xs">Demons</span>
					<span class="bg-pink-500 px-2 py-1 rounded text-xs">Shounen</span>
				</div>
			</div>
		</div>
	</div>
	
	<!-- Tabs -->
	<div class="border-b border-gray-200 mb-8">
		<nav class="flex space-x-8">
			<button class="tab-btn active py-3 px-1 border-b-2 border-primary text-primary font-semibold" 
			        onclick="showTab('overview')">Overview</button>
			<button class="tab-btn py-3 px-1 border-b-2 border-transparent text-gray-500 hover:text-gray-700 font-medium" 
			        onclick="showTab('episodes')">Episodes</button>
			<button class="tab-btn py-3 px-1 border-b-2 border-transparent text-gray-500 hover:text-gray-700 font-medium" 
			        onclick="showTab('comments')">Comments</button>
		</nav>
	</div>
	
	<!-- Tab Content -->
	<div id="tab-content">
		<!-- Overview Tab -->
		<div id="overview-tab" class="tab-content">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
				<div class="space-y-4">
					<div class="flex justify-between"><span class="text-gray-500 font-medium">Type:</span> <span class="text-gray-800 font-semibold">%s</span></div>
					<div class="flex justify-between"><span class="text-gray-500 font-medium">Status:</span> <span class="text-gray-800 font-semibold capitalize">%s</span></div>
					<div class="flex justify-between"><span class="text-gray-500 font-medium">Episodes:</span> <span class="text-gray-800 font-semibold">%s</span></div>
				</div>
				<div class="space-y-4">
					<div class="flex justify-between"><span class="text-gray-500 font-medium">Year:</span> <span class="text-gray-800 font-semibold">%s</span></div>
					<div class="flex justify-between"><span class="text-gray-500 font-medium">Genres:</span> <span class="text-gray-800 font-semibold">%s</span></div>
				</div>
			</div>
			<div class="mb-8">
				<h3 class="text-2xl font-bold text-gray-800 mb-4">Synopsis</h3>
				<p class="text-gray-600 leading-relaxed text-lg">%s</p>
			</div>
			<div id="anime-actions" class="flex flex-wrap gap-3">
				<!-- Actions will be populated by JavaScript based on auth status -->
			</div>
			<script>
				// Populate actions based on auth status
				const actionsDiv = document.getElementById('anime-actions');
				const isLoggedIn = localStorage.getItem('auth_token');
				
				if (isLoggedIn) {
					actionsDiv.innerHTML = '<button class="bg-gradient-to-r from-green-500 to-emerald-600 hover:from-green-600 hover:to-emerald-700 text-white px-6 py-3 rounded-xl font-semibold transition-all" onclick="addAnimeToList(\'%s\', \'watching\')">Add to Watching</button>' +
						'<button class="bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white px-6 py-3 rounded-xl font-semibold transition-all" onclick="addAnimeToList(\'%s\', \'completed\')">Mark Completed</button>' +
						'<button class="bg-gradient-to-r from-yellow-500 to-orange-600 hover:from-yellow-600 hover:to-orange-700 text-white px-6 py-3 rounded-xl font-semibold transition-all" onclick="addAnimeToList(\'%s\', \'plan-to-watch\')">Plan to Watch</button>' +
						'<button class="bg-gradient-to-r from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700 text-white px-6 py-3 rounded-xl font-semibold transition-all" onclick="addAnimeToList(\'%s\', \'on-hold\')">On Hold</button>' +
						'<button class="bg-gradient-to-r from-red-500 to-pink-600 hover:from-red-600 hover:to-pink-700 text-white px-6 py-3 rounded-xl font-semibold transition-all" onclick="addAnimeToList(\'%s\', \'dropped\')">Dropped</button>';
				} else {
					actionsDiv.innerHTML = '<button class="bg-gradient-to-r from-purple-500 to-indigo-600 hover:from-purple-600 hover:to-indigo-700 text-white px-6 py-3 rounded-xl font-semibold transition-all" onclick="watchTrailer()">Watch Trailer</button>' +
						'<button class="bg-gradient-to-r from-gray-500 to-gray-600 hover:from-gray-600 hover:to-gray-700 text-white px-6 py-3 rounded-xl font-semibold transition-all" onclick="signIn()">Sign In to Add to List</button>';
				}
				
				function watchTrailer() {
					alert('Trailer functionality coming soon!');
				}
			</script>
		</div>
		
		<!-- Episodes Tab -->
		<div id="episodes-tab" class="tab-content hidden">
			<div class="mb-4">
				<button class="bg-anime-blue hover:bg-blue-600 px-4 py-2 rounded-lg mr-2" onclick="loadEpisodes('%s', 1)">Season 1</button>
				<button class="bg-gray-600 hover:bg-gray-500 px-4 py-2 rounded-lg" onclick="loadEpisodes('%s', 2)">Season 2</button>
			</div>
			<div id="episode-list" class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<!-- Episodes will be loaded here -->
				<div class="bg-gray-700 rounded-lg p-4 hover:bg-gray-600 cursor-pointer" onclick="playEpisode(1)">
					<div class="flex items-center space-x-3">
						<div class="w-16 h-12 bg-gray-600 rounded flex items-center justify-center">
							<svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
								<path d="M8 5v10l8-5-8-5z"/>
							</svg>
						</div>
						<div>
							<h4 class="text-white font-medium">Episode 1</h4>
							<p class="text-gray-400 text-sm">24 min</p>
						</div>
					</div>
				</div>
			</div>
		</div>
		
		<!-- Comments Tab -->
		<div id="comments-tab" class="tab-content hidden">
			<form class="mb-6" hx-post="/api/anime/%s/comments" hx-target="#comments-list" hx-swap="afterbegin">
				<textarea name="content" placeholder="Write a comment..." class="w-full bg-gray-700 border border-gray-600 rounded-lg p-3 text-white placeholder-gray-400 resize-none" rows="3"></textarea>
				<div class="flex justify-between items-center mt-2">
					<label class="flex items-center text-gray-400">
						<input type="checkbox" name="spoiler" class="mr-2"> Contains spoilers
					</label>
					<button type="submit" class="bg-anime-blue hover:bg-blue-600 px-4 py-2 rounded-lg transition-colors">Post Comment</button>
				</div>
			</form>
			<div id="comments-list" hx-get="/api/anime/%s/comments" hx-trigger="load">
				<!-- Comments will be loaded here -->
			</div>
		</div>
	</div>
	
	<script>
		function showTab(tabName) {
			// Hide all tabs
			document.querySelectorAll('.tab-content').forEach(tab => tab.classList.add('hidden'));
			document.querySelectorAll('.tab-btn').forEach(btn => {
				btn.classList.remove('border-anime-blue', 'text-anime-blue');
				btn.classList.add('border-transparent', 'text-gray-400');
			});
			
			// Show selected tab
			document.getElementById(tabName + '-tab').classList.remove('hidden');
			event.target.classList.remove('border-transparent', 'text-gray-400');
			event.target.classList.add('border-anime-blue', 'text-anime-blue');
		}
		
		function loadEpisodes(animeId, season) {
			// Mock episode loading
			console.log('Loading episodes for anime:', animeId, 'season:', season);
		}
		
		function playEpisode(episodeNumber) {
			document.getElementById('episode-title').textContent = 'Episode ' + episodeNumber;
			document.getElementById('video-player').innerHTML = '<div class="flex items-center justify-center h-full text-white"><p>Video player would load here</p></div>';
			showPlayer();
		}
	</script>`,
		anime.Name, anime.BannerUrl, anime.Name, anime.ImageUrl, anime.Name, 
		anime.Score, anime.Type, yearSeason, progressText,
		anime.Type, anime.Status, progressText, yearSeason, genreStr, anime.Notes,
		anime.Name, anime.Name, anime.Name, anime.Name, anime.Name,
		anime.ID.Hex(), anime.ID.Hex(),
		anime.ID.Hex(), anime.ID.Hex())
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
	
	// Get user ID if authenticated
	userID := ""
	if user := r.Context().Value("user"); user != nil {
		claims := user.(*middleware.SupabaseClaims)
		userID = claims.Sub
	}

	filteredAnimes := services.SmartSearch(search, genre, year, season, format, status, userID)
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
	http.ServeFile(w, r, "./static/index.html")
}

func ServeOldFrontendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html lang="en" class="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AnimeVerse - Modern Anime Database</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/@supabase/supabase-js@2"></script>
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
    </style>
</head>
<body class="bg-gray-900 text-white min-h-screen">
    <header class="bg-gray-800 shadow-lg">
        <div class="container mx-auto px-4 py-6">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-4">
                    <div class="text-3xl font-bold bg-gradient-to-r from-anime-blue to-anime-purple bg-clip-text text-transparent">
                        🌸 AnimeVerse
                    </div>
                    <span class="text-sm text-gray-400">v3.0</span>
                </div>
                <nav class="hidden md:flex space-x-4 items-center">
                    <a href="#" class="text-gray-300 hover:text-white transition-colors">Home</a>
                    <a href="/api/animes" class="text-gray-300 hover:text-white transition-colors">API</a>
                    <div id="auth-section">
                        <div id="signed-out" class="flex space-x-2">
                            <button onclick="signIn()" class="bg-anime-blue hover:bg-blue-600 px-3 py-1 rounded text-sm transition-colors">Sign In</button>
                            <button onclick="signUp()" class="bg-anime-purple hover:bg-purple-600 px-3 py-1 rounded text-sm transition-colors">Sign Up</button>
                        </div>
                        <div id="signed-in" class="hidden flex space-x-2 items-center">
                            <span id="user-name" class="text-gray-300"></span>
                            <button onclick="signOut()" class="bg-gray-600 hover:bg-gray-500 px-3 py-1 rounded text-sm transition-colors">Sign Out</button>
                        </div>
                    </div>
                    <div id="user-controls" class="hidden flex space-x-2">
                        <button class="bg-green-600 hover:bg-green-700 px-3 py-1 rounded text-sm transition-colors" onclick="searchExternalAnime()">Add Anime</button>
                    </div>
                    <div id="admin-controls" class="hidden flex space-x-2">
                        <button class="bg-anime-blue hover:bg-blue-600 px-3 py-1 rounded text-sm transition-colors" onclick="importTrending()">Import Trending</button>
                        <button class="bg-anime-purple hover:bg-purple-600 px-3 py-1 rounded text-sm transition-colors" onclick="importSeasonal()">Import Seasonal</button>
                    </div>
                </nav>
            </div>
        </div>
    </header>
    <section class="bg-gray-800 border-b border-gray-700">
        <div class="container mx-auto px-4 py-6">
            <form id="filter-form" class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4">
                <input type="text" name="search" id="search-input" placeholder="Search anime..." class="bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-anime-blue">
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
    <main class="container mx-auto px-4 py-8">
        <section>
            <h2 class="text-2xl font-bold mb-6">📚 All Anime</h2>
            <div id="anime-grid" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
                <div class="animate-pulse bg-gray-700 rounded-lg h-80"></div>
                <div class="animate-pulse bg-gray-700 rounded-lg h-80"></div>
                <div class="animate-pulse bg-gray-700 rounded-lg h-80"></div>
            </div>
        </section>
    </main>
    <div id="modal" class="fixed inset-0 bg-black bg-opacity-50 hidden items-center justify-center z-50">
        <div class="bg-gray-800 rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
            <div id="modal-content"></div>
            <button onclick="closeModal()" class="mt-4 bg-anime-blue hover:bg-blue-600 px-4 py-2 rounded-lg transition-colors">Close</button>
        </div>
    </div>
    <script>
        const supabaseUrl = 'https://rrrgpcnhzmnnjacvgzcn.supabase.co';
        const supabaseKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJycmdwY25oem1ubmphY3ZnemNuIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTA3NDY3NDcsImV4cCI6MjA2NjMyMjc0N30.zZQ2aP_jEbCuqG3JehTni3xAesXrafiHUFYeD_-tTcE';
        const supabase = window.supabase.createClient(supabaseUrl, supabaseKey);
        let currentUser = null;
        let authToken = null;
        function showModal() { document.getElementById('modal').classList.remove('hidden'); document.getElementById('modal').classList.add('flex'); }
        function closeModal() { document.getElementById('modal').classList.add('hidden'); document.getElementById('modal').classList.remove('flex'); }
        function signIn() { showAuthModal('login'); }
        function signUp() { showAuthModal('register'); }
        function signOut() { localStorage.removeItem('auth_token'); authToken = null; currentUser = null; updateAuthUI(); applyFilters(); }
        function showAuthModal(type) {
            const title = type === 'login' ? 'Sign In' : 'Sign Up';
            const buttonText = type === 'login' ? 'Sign In' : 'Sign Up';
            const switchText = type === 'login' ? 'Need an account? Sign up' : 'Have an account? Sign in';
            const switchAction = type === 'login' ? 'register' : 'login';
            document.getElementById('modal-content').innerHTML = '<div class="max-w-md mx-auto"><h2 class="text-2xl font-bold text-white mb-6">' + title + '</h2><form id="auth-form" class="space-y-4">' + (type === 'register' ? '<input type="text" name="name" placeholder="Full Name" class="w-full bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-anime-blue">' : '') + '<input type="email" name="email" placeholder="Email" required class="w-full bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-anime-blue"><input type="password" name="password" placeholder="Password" required class="w-full bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-anime-blue"><button type="submit" class="w-full bg-anime-blue hover:bg-blue-600 px-4 py-2 rounded-lg transition-colors">' + buttonText + '</button></form><div class="mt-4 text-center"><button onclick="showAuthModal(\'' + switchAction + '\')" class="text-anime-blue hover:text-blue-400 text-sm">' + switchText + '</button></div><div class="mt-4 text-center"><p class="text-gray-400 text-sm mb-2">Or continue with:</p><div class="flex gap-2 justify-center"><button onclick="oauthLogin(\'google\')" class="bg-red-600 hover:bg-red-700 px-4 py-2 rounded text-sm transition-colors">Google</button><button onclick="oauthLogin(\'github\')" class="bg-gray-700 hover:bg-gray-600 px-4 py-2 rounded text-sm transition-colors">GitHub</button></div></div><div class="mt-4 p-3 bg-gray-700 rounded text-sm text-gray-300"><strong>Demo Account:</strong><br>Email: demo@animeverse.com<br>Password: demo123</div></div>';
            document.getElementById('auth-form').addEventListener('submit', function(e) { e.preventDefault(); handleAuth(type); });
            showModal();
        }
        function handleAuth(type) {
            const form = document.getElementById('auth-form');
            const formData = new FormData(form);
            const data = { email: formData.get('email'), password: formData.get('password'), name: formData.get('name') || '' };
            fetch('/auth/' + type, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    authToken = data.token;
                    currentUser = data.user;
                    localStorage.setItem('auth_token', authToken);
                    updateAuthUI();
                    closeModal();
                    applyFilters();
                } else {
                    alert(data.message || 'Authentication failed');
                }
            })
            .catch(error => { console.error('Auth error:', error); alert('Authentication failed'); });
        }
        function oauthLogin(provider) { supabase.auth.signInWithOAuth({ provider: provider, options: { redirectTo: window.location.origin } }); }
        function updateAuthUI() {
            const signedOut = document.getElementById('signed-out');
            const signedIn = document.getElementById('signed-in');
            const userControls = document.getElementById('user-controls');
            const adminControls = document.getElementById('admin-controls');
            const userName = document.getElementById('user-name');
            if (currentUser) {
                signedOut.classList.add('hidden');
                signedIn.classList.remove('hidden');
                userControls.classList.remove('hidden');
                userName.textContent = currentUser.name || currentUser.email;
                if (currentUser.role === 'admin') { adminControls.classList.remove('hidden'); }
            } else {
                signedOut.classList.remove('hidden');
                signedIn.classList.add('hidden');
                userControls.classList.add('hidden');
                adminControls.classList.add('hidden');
            }
        }
        function importTrending() {
            if (!authToken) return;
            fetch('/api/admin/import/trending', { method: 'POST', headers: { 'Authorization': 'Bearer ' + authToken, 'Content-Type': 'application/json' } })
            .then(response => response.json())
            .then(data => { alert(data.message || 'Import completed'); applyFilters(); })
            .catch(error => { console.error('Import error:', error); alert('Import failed'); });
        }
        function importSeasonal() {
            if (!authToken) return;
            fetch('/api/admin/import/seasonal?year=2024&season=winter', { method: 'POST', headers: { 'Authorization': 'Bearer ' + authToken, 'Content-Type': 'application/json' } })
            .then(response => response.json())
            .then(data => { alert(data.message || 'Import completed'); applyFilters(); })
            .catch(error => { console.error('Import error:', error); alert('Import failed'); });
        }
        function addAnimeToList(animeName, status) {
            if (!authToken) { alert('Please sign in first'); return; }
            fetch('/api/user/anime', {
                method: 'POST',
                headers: { 'Authorization': 'Bearer ' + authToken, 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: animeName, status: status })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Added to ' + status + ' list!');
                    applyFilters();
                } else {
                    alert(data.error || 'Failed to add anime');
                }
            })
            .catch(error => { console.error('Add error:', error); alert('Failed to add anime'); });
        }
        function updateAnimeStatus(animeId, status) {
            if (!authToken) return;
            fetch('/api/user/anime/' + animeId + '/status', {
                method: 'PUT',
                headers: { 'Authorization': 'Bearer ' + authToken, 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: status })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Status updated to ' + status);
                    applyFilters();
                }
            })
            .catch(error => console.error('Update error:', error));
        }
        function searchExternalAnime() {
            if (!authToken) { alert('Please sign in first'); return; }
            const query = prompt('Search for anime:');
            if (!query) return;
            fetch('/api/user/search?q=' + encodeURIComponent(query), {
                headers: { 'Authorization': 'Bearer ' + authToken }
            })
            .then(response => response.json())
            .then(data => {
                if (data.success && data.data.length > 0) {
                    let options = 'Found anime:\n';
                    data.data.forEach((anime, i) => {
                        options += (i + 1) + '. ' + anime.name + ' (' + anime.year + ')\n';
                    });
                    const choice = prompt(options + '\nEnter number to add to Plan to Watch:');
                    if (choice && data.data[choice - 1]) {
                        addAnimeToList(data.data[choice - 1].name, 'plan-to-watch');
                    }
                } else {
                    alert('No anime found');
                }
            })
            .catch(error => { console.error('Search error:', error); alert('Search failed'); });
        }
        function applyFilters() {
            const form = document.getElementById('filter-form');
            const formData = new FormData(form);
            const params = new URLSearchParams();
            for (let [key, value] of formData.entries()) {
                if (value.trim() !== '') { params.append(key, value); }
            }
            document.getElementById('loading').classList.remove('hidden');
            document.getElementById('anime-grid').style.opacity = '0.5';
            const headers = { 'HX-Request': 'true' };
            if (authToken) { headers['Authorization'] = 'Bearer ' + authToken; }
            fetch('/api/animes/filter?' + params.toString(), { headers })
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
        let searchTimeout;
        function debounceSearch() { clearTimeout(searchTimeout); searchTimeout = setTimeout(applyFilters, 500); }
        document.addEventListener('DOMContentLoaded', async function() {
            authToken = localStorage.getItem('auth_token');
            if (authToken) {
                fetch('/api/user/me', { headers: { 'Authorization': 'Bearer ' + authToken } })
                .then(response => { if (response.ok) { return response.json(); } throw new Error('Invalid token'); })
                .then(data => { if (data.success) { currentUser = data.data; } updateAuthUI(); })
                .catch(() => { localStorage.removeItem('auth_token'); authToken = null; updateAuthUI(); });
            } else {
                const { data: { session } } = await supabase.auth.getSession();
                if (session) {
                    authToken = session.access_token;
                    currentUser = { id: session.user.id, email: session.user.email, name: session.user.user_metadata.full_name || session.user.email };
                    localStorage.setItem('auth_token', authToken);
                }
                updateAuthUI();
            }
            supabase.auth.onAuthStateChange((event, session) => {
                if (event === 'SIGNED_IN' && session) {
                    authToken = session.access_token;
                    currentUser = { id: session.user.id, email: session.user.email, name: session.user.user_metadata.full_name || session.user.email };
                    localStorage.setItem('auth_token', authToken);
                    updateAuthUI();
                    window.location.reload();
                } else if (event === 'SIGNED_OUT') {
                    localStorage.removeItem('auth_token');
                    authToken = null;
                    currentUser = null;
                    updateAuthUI();
                }
            });
            document.getElementById('search-input').addEventListener('input', debounceSearch);
            document.getElementById('genre-select').addEventListener('change', applyFilters);
            document.getElementById('year-select').addEventListener('change', applyFilters);
            document.getElementById('season-select').addEventListener('change', applyFilters);
            document.getElementById('format-select').addEventListener('change', applyFilters);
            document.getElementById('status-select').addEventListener('change', applyFilters);
            applyFilters();
        });
    </script>
</body>
</html>`))
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

func BackfillDataHandler(w http.ResponseWriter, r *http.Request) {
	count, err := services.BackfillAllMissingData()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to backfill data: "+err.Error())
		return
	}
	sendJSONResponse(w, http.StatusOK, true, fmt.Sprintf("Backfilled %d anime records", count), nil, "")
}

func BulkImportHandler(w http.ResponseWriter, r *http.Request) {
	count, err := services.BulkImportAnimeDatabase()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to bulk import: "+err.Error())
		return
	}
	sendJSONResponse(w, http.StatusOK, true, fmt.Sprintf("Successfully imported %d anime from database", count), nil, "")
}

func UpdateCurrentSeasonHandler(w http.ResponseWriter, r *http.Request) {
	count, err := services.UpdateCurrentSeasonAnime()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to update current season: "+err.Error())
		return
	}
	sendJSONResponse(w, http.StatusOK, true, fmt.Sprintf("Updated %d current season anime", count), nil, "")
}

func SearchAnimesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Search query is required")
		return
	}
	
	searchResults := services.SearchAnimes(query)
	if searchResults == nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to search animes")
		return
	}
	
	// Check if request wants HTML (HTMX)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderSearchResults(w, searchResults)
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Search results retrieved successfully", searchResults, "")
}

// renderSearchResults renders search results as cards
func renderSearchResults(w http.ResponseWriter, animes []primitive.M) {
	if len(animes) == 0 {
		fmt.Fprintf(w, `<div class="col-span-full text-center py-12 text-gray-500">
			<p class="text-xl mb-2">😢 No anime found</p>
			<p>Try a different search term</p>
		</div>`)
		return
	}
	
	for _, anime := range animes {
		name := ""
		imageUrl := "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjI2NyIgdmlld0JveD0iMCAwIDIwMCAyNjciIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjY3IiBmaWxsPSIjRjNGNEY2Ii8+Cjx0ZXh0IHg9IjEwMCIgeT0iMTMzIiBmaWxsPSIjOUI5QkE0IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNCI+Tm8gSW1hZ2U8L3RleHQ+Cjwvc3ZnPg=="
		score := 0
		year := 0
		episodes := "N/A"
		
		if n, ok := anime["name"].(string); ok {
			name = n
		}
		if img, ok := anime["imageUrl"].(string); ok && img != "" {
			imageUrl = img
		}
		if s, ok := anime["score"].(int32); ok {
			score = int(s)
		} else if s, ok := anime["score"].(int); ok {
			score = s
		}
		if y, ok := anime["year"].(int32); ok {
			year = int(y)
		} else if y, ok := anime["year"].(int); ok {
			year = y
		}
		if prog, ok := anime["progress"].(primitive.M); ok {
			if total, ok := prog["total"].(int32); ok && total > 0 {
				episodes = fmt.Sprintf("%d Episodes", total)
			}
		}
		
		fmt.Fprintf(w, `
		<div class="bg-white rounded-2xl shadow-lg overflow-hidden anime-card cursor-pointer"
		     onclick="closeModal(); setTimeout(() => showAnimeModal('%s'), 100);">
		    <div class="aspect-[3/4] bg-gray-200 relative overflow-hidden">
		        <img src="%s" alt="%s" class="w-full h-full object-cover">
		        <div class="absolute top-3 right-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white px-2 py-1 rounded-full text-xs font-bold">
		            %d★
		        </div>
		    </div>
		    <div class="p-4">
		        <h3 class="font-bold text-gray-800 text-sm mb-2 truncate" title="%s">%s</h3>
		        <div class="flex items-center justify-between text-xs text-gray-500">
		            <span>%d</span>
		            <span class="bg-gray-100 px-2 py-1 rounded-full">%s</span>
		        </div>
		    </div>
		</div>`, 
			strings.ReplaceAll(strings.ToLower(name), " ", "-"), 
			imageUrl, name, score, name, name, year, episodes)
	}
}

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "", nil, "Not authenticated")
		return
	}
	
	claims := user.(*middleware.SupabaseClaims)
	
	// Get or create user in MongoDB
	dbUser, err := services.CreateOrUpdateUser(claims.Sub, claims.Email, claims.Name)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to get user data")
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "User retrieved successfully", dbUser, "")
}

func GetUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "", nil, "Not authenticated")
		return
	}
	
	claims := user.(*middleware.SupabaseClaims)
	
	// Update stats first
	if err := services.UpdateUserStats(claims.Sub); err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to update user stats")
		return
	}
	
	stats, err := services.GetUserStats(claims.Sub)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to get user stats")
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "User stats retrieved successfully", stats, "")
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

// Image caching handlers
func CheckImagesHandler(w http.ResponseWriter, r *http.Request) {
	malIDStr := r.URL.Query().Get("mal_id")
	anilistIDStr := r.URL.Query().Get("anilist_id")
	
	malID := 0
	anilistID := 0
	
	if malIDStr != "" {
		if id, err := strconv.Atoi(malIDStr); err == nil {
			malID = id
		}
	}
	
	if anilistIDStr != "" {
		if id, err := strconv.Atoi(anilistIDStr); err == nil {
			anilistID = id
		}
	}
	
	images := services.GetImagesByIDs(malID, anilistID)
	if images != nil {
		sendJSONResponse(w, http.StatusOK, true, "Images found", images, "")
	} else {
		sendJSONResponse(w, http.StatusNotFound, false, "Images not found", nil, "")
	}
}

func SaveImagesHandler(w http.ResponseWriter, r *http.Request) {
	var req model.ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}
	
	if err := services.SaveImageData(req); err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to save images")
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Images saved successfully", nil, "")
}

func GetRandomAnimeHandler(w http.ResponseWriter, r *http.Request) {
	randomAnime := services.GetRandomAnime()
	if randomAnime == nil {
		sendJSONResponse(w, http.StatusNotFound, false, "", nil, "No anime found")
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Random anime retrieved", randomAnime, "")
}

func GetTop2025AnimesHandler(w http.ResponseWriter, r *http.Request) {
	top2025Animes := services.GetTop2025Animes()
	if top2025Animes == nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to fetch top 2025 animes")
		return
	}
	
	// Check if request wants HTML (HTMX)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderHeroCards(w, top2025Animes)
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Top 2025 animes retrieved successfully", top2025Animes, "")
}

func GetPreviewAnimesHandler(w http.ResponseWriter, r *http.Request) {
	previewAnimes := services.GetPreviewAnimes()
	if previewAnimes == nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to fetch preview animes")
		return
	}
	
	// Check if request wants HTML (HTMX)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderAnimeCards(w, previewAnimes)
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Preview animes retrieved successfully", previewAnimes, "")
}

// renderHeroCards renders anime cards for hero section
func renderHeroCards(w http.ResponseWriter, animes []primitive.M) {
	if len(animes) == 0 {
		fmt.Fprintf(w, `<div class="flex-none w-48 bg-white/10 backdrop-blur rounded-2xl p-4">
			<div class="aspect-[3/4] bg-white/20 rounded-xl mb-3 flex items-center justify-center">
				<span class="text-white/50 text-sm">No anime found</span>
			</div>
		</div>`)
		return
	}
	
	for _, anime := range animes {
		name := ""
		imageUrl := "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjI2NyIgdmlld0JveD0iMCAwIDIwMCAyNjciIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iMjY3IiBmaWxsPSIjRjNGNEY2Ii8+Cjx0ZXh0IHg9IjEwMCIgeT0iMTMzIiBmaWxsPSIjOUI5QkE0IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNCI+Tm8gSW1hZ2U8L3RleHQ+Cjwvc3ZnPg=="
		episodes := "N/A"
		
		if n, ok := anime["name"].(string); ok {
			name = n
		}
		if img, ok := anime["imageUrl"].(string); ok && img != "" {
			imageUrl = img
		}
		if prog, ok := anime["progress"].(primitive.M); ok {
			if total, ok := prog["total"].(int32); ok && total > 0 {
				episodes = fmt.Sprintf("%d Episodes", total)
			}
		}
		
		fmt.Fprintf(w, `
		<div class="flex-none w-48 bg-white/10 backdrop-blur rounded-2xl p-4 hover:bg-white/20 transition-all cursor-pointer"
		     onclick="showAnimeModal('%s')">
		    <div class="aspect-[3/4] bg-white/20 rounded-xl mb-3 overflow-hidden">
		        <img src="%s" alt="%s" class="w-full h-full object-cover">
		    </div>
		    <h3 class="font-bold text-sm mb-1 text-white truncate">%s</h3>
		    <p class="text-white/70 text-xs">%s</p>
		</div>`, 
			strings.ReplaceAll(strings.ToLower(name), " ", "-"), 
			imageUrl, name, name, episodes)
	}
}

func GetScheduleHandler(w http.ResponseWriter, r *http.Request) {
	day := r.URL.Query().Get("day")
	if day == "" {
		day = time.Now().Weekday().String()
	}
	
	schedule := services.GetDailySchedule(day)
	
	// Check if request wants HTML (HTMX)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		renderSchedule(w, schedule)
		return
	}
	
	sendJSONResponse(w, http.StatusOK, true, "Schedule retrieved", schedule, "")
}

func renderSchedule(w http.ResponseWriter, schedule []map[string]interface{}) {
	if len(schedule) == 0 {
		fmt.Fprintf(w, `<div class="text-center py-8 text-gray-500">
			<p>No scheduled anime for today</p>
		</div>`)
		return
	}
	
	for _, item := range schedule {
		name := ""
		time := ""
		episode := ""
		
		if n, ok := item["name"].(string); ok {
			name = n
		}
		if t, ok := item["time"].(string); ok {
			time = t
		}
		if e, ok := item["episode"].(string); ok {
			episode = e
		}
		
		fmt.Fprintf(w, `
		<div class="flex items-center justify-between p-4 hover:bg-gray-50 rounded-xl transition-colors cursor-pointer border-l-4 border-primary"
		     onclick="showAnimeModal('%s')">
		    <div class="flex-1">
		        <h3 class="font-semibold text-gray-800 mb-1">%s</h3>
		        <p class="text-gray-500 text-sm">%s</p>
		    </div>
		    <div class="text-right">
		        <span class="bg-gradient-to-r from-indigo-500 to-purple-600 text-white px-3 py-1 rounded-full text-sm font-medium">%s</span>
		    </div>
		</div>`,
			strings.ReplaceAll(strings.ToLower(name), " ", "-"),
			name, episode, time)
	}
}
