package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Flack74/mongoapi/config"
	model "github.com/Flack74/mongoapi/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InsertOneAnime inserts a new anime document into MongoDB.
func InsertOneAnime(anime model.Anime) {
	inserted, err := config.Collection.InsertOne(context.Background(), anime)
	if err != nil {
		log.Println("Error inserting anime:", err)
		return
	}
	fmt.Println("Inserted 1 anime in db with id:", inserted.InsertedID)
}

// UpdateAnime performs a generic partial update on any fields provided in the JSON body.
func UpdateAnime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idStr := params["id"]

	// Convert the string ID to an ObjectID.
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Decode request body into a map.
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Remove _id if present to prevent overwriting the document's ID.
	delete(updates, "_id")

	// Construct the update query using $set.
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updates}

	result, err := config.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating anime:", err)
		http.Error(w, "Failed to update anime", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		log.Println("No document matched the given ID.")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Println("Successfully updated anime:", idStr)
	// Respond with the update result.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"matched":  result.MatchedCount,
		"modified": result.ModifiedCount,
		"id":       idStr,
	})
}

// DeleteOneAnime deletes a single anime by ID.
func DeleteOneAnime(animeId string) {
	id, err := primitive.ObjectIDFromHex(animeId)
	if err != nil {
		log.Println("Invalid ObjectID:", err)
		return
	}

	filter := bson.M{"_id": id}
	result, err := config.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Error deleting anime:", err)
		return
	}
	fmt.Println("Anime deleted with count:", result.DeletedCount)
}

// DeleteAllAnime deletes all anime records from the collection.
func DeleteAllAnime() int64 {
	deleteResult, err := config.Collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Println("Error deleting all animes:", err)
		return 0
	}
	fmt.Println("Number of animes deleted:", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

// GetAllAnimes retrieves all anime documents.
func GetAllAnimes() []primitive.M {
	cur, err := config.Collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Println("Error fetching animes:", err)
		return nil
	}
	defer cur.Close(context.Background())

	var animes []primitive.M
	for cur.Next(context.Background()) {
		var anime bson.M
		if err := cur.Decode(&anime); err != nil {
			log.Println("Error decoding anime:", err)
			continue
		}
		animes = append(animes, anime)
	}
	return animes
}

// Controller Handlers

// GetMyAllAnimes handles GET /api/animes.
func GetMyAllAnimes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allAnimes := GetAllAnimes()
	json.NewEncoder(w).Encode(allAnimes)
}

// CreateAnime handles POST /api/anime.
func CreateAnime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var anime model.Anime
	if err := json.NewDecoder(r.Body).Decode(&anime); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	InsertOneAnime(anime)
	json.NewEncoder(w).Encode(anime)
}

// UpdateAnime handler is used for PUT /api/anime/{id}.
// It updates the document with any fields provided in the request body.
func UpdateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	UpdateAnime(w, r)
}

// DeleteAnAnime handles DELETE /api/anime/{id}.
func DeleteAnAnime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	DeleteOneAnime(params["id"])
	json.NewEncoder(w).Encode(map[string]string{"deleted": params["id"]})
}

// DeleteEveryAnimes handles DELETE /api/deleteallanime.
func DeleteEveryAnimes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	count := DeleteAllAnime()
	json.NewEncoder(w).Encode(map[string]int64{"deleted_count": count})
}

// ServeHome serves the home page.
func ServeHome(w http.ResponseWriter, r *http.Request) {
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
