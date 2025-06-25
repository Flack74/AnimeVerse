package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Flack74/mongoapi/middleware"
	model "github.com/Flack74/mongoapi/models"
	"github.com/Flack74/mongoapi/services"
	"github.com/go-chi/chi/v5"
)

func AddAnimeHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "", nil, "Not authenticated")
		return
	}

	claims := user.(*middleware.SupabaseClaims)

	var req struct {
		Name   string             `json:"name"`
		Status model.WatchStatus  `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request")
		return
	}

	anime, err := services.AddAnimeToUserList(claims.Sub, req.Name, req.Status)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to add anime")
		return
	}

	sendJSONResponse(w, http.StatusCreated, true, "Anime added successfully", anime, "")
}

func UpdateAnimeStatusHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "", nil, "Not authenticated")
		return
	}

	claims := user.(*middleware.SupabaseClaims)
	animeID := chi.URLParam(r, "id")

	var req struct {
		Status model.WatchStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request")
		return
	}

	anime, err := services.UpdateAnimeStatus(claims.Sub, animeID, req.Status)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to update status")
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "Status updated successfully", anime, "")
}

func UpdateAnimeScoreHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "", nil, "Not authenticated")
		return
	}

	claims := user.(*middleware.SupabaseClaims)
	animeID := chi.URLParam(r, "id")

	var req struct {
		Score int `json:"score"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request")
		return
	}

	anime, err := services.UpdateAnimeScore(claims.Sub, animeID, req.Score)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to update score")
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "Score updated successfully", anime, "")
}

func RemoveAnimeHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "", nil, "Not authenticated")
		return
	}

	claims := user.(*middleware.SupabaseClaims)
	animeID := chi.URLParam(r, "id")

	err := services.RemoveAnimeFromUserList(claims.Sub, animeID)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to remove anime")
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "Anime removed successfully", nil, "")
}

func SearchAnimeHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		sendJSONResponse(w, http.StatusUnauthorized, false, "", nil, "Not authenticated")
		return
	}

	claims := user.(*middleware.SupabaseClaims)
	query := r.URL.Query().Get("q")

	if query == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "", nil, "Query parameter required")
		return
	}

	results, err := services.SearchAndAddAnime(claims.Sub, query)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", nil, "Search failed")
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "Search completed", results, "")
}