package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	model "github.com/Flack74/mongoapi/models"
	"github.com/Flack74/mongoapi/services"
	"github.com/gorilla/mux"
)


func GetMyAllAnimesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allAnimes := services.GetAllAnimes()
	json.NewEncoder(w).Encode(allAnimes)
}

func GetAnimeByNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	animeName := vars["animeName"]

	animeName = strings.ReplaceAll(animeName, "-", " ")
	animeName = strings.ReplaceAll(animeName, "_", " ")

	animeName = strings.ToLower(animeName)

	existingAnime, err := services.SearchAnimeByName(animeName)
	if err != nil || existingAnime == nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(existingAnime)
}

func CreateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var anime model.Anime
	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingAnime, err := services.FindAnimeByName(anime.Name)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}

	if existingAnime != nil {
		http.Error(w, "Anime with this name already exists", http.StatusConflict)
		return
	}

	err = services.InsertOneAnime(anime)

	if err != nil {
		http.Error(w, "Failed to insert anime", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(anime)
}

func UpdateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	services.UpdateAnime(w, r)
}

func DeleteAnAnimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	services.DeleteOneAnime(params["id"])
	json.NewEncoder(w).Encode(map[string]string{"deleted": params["id"]})
}

func DeleteEveryAnimesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	count := services.DeleteAllAnime()
	json.NewEncoder(w).Encode(map[string]int64{"deleted_count": count})
}

func ServeHomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Anime API by Flack</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #1e1e1e;
					color: #fff;
					text-align: center;
					padding: 50px;
				}
				h1 {
					font-size: 2.5rem;
					color: #ffcc00;
				}
				.container {
					max-width: 600px;
					margin: auto;
					padding: 20px;
					background: #2a2a2a;
					border-radius: 10px;
					box-shadow: 0 4px 10px rgba(255, 204, 0, 0.3);
				}
				p {
					font-size: 1.2rem;
					margin-top: 10px;
					color: #ddd;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Welcome to Anime API by Flack</h1>
				<p>Explore a collection of anime data through our simple API.</p>
			</div>
		</body>
		</html>
	`))
}
