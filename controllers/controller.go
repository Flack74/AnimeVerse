package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	model "github.com/Flack74/mongoapi/models"
	"github.com/Flack74/mongoapi/services"
	"github.com/go-chi/chi/v5"
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

func GetMyAllAnimesHandler(w http.ResponseWriter, r *http.Request) {
	allAnimes := services.GetAllAnimes()
	if allAnimes == nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to fetch animes")
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
		"version":   "2.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}, "")
}

func ServeHomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>AnimeVerse API v2.0</title>
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
				<div class="version">v2.0 - Powered by Chi Router</div>
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
