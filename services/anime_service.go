package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Flack74/mongoapi/config"
	model "github.com/Flack74/mongoapi/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAnimeByName(name string) (*model.Anime, error) {
	collection := config.Collection // Ensure this matches your database collection name
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var anime model.Anime
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&anime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No existing anime found
		}
		log.Println("Error finding anime:", err)
		return nil, err
	}
	return &anime, nil
}

func SearchAnimeByName(name string) (*model.Anime, error) {
	collection := config.Collection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var anime model.Anime
	filter := bson.M{"name": bson.M{"$regex": "^" + name + "$", "$options": "i"}} // Case-insensitive search

	err := collection.FindOne(ctx, filter).Decode(&anime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Println("Error finding anime:", err)
		return nil, err
	}
	return &anime, nil
}

// InsertOneAnime inserts a new anime document into MongoDB.
func InsertOneAnime(anime model.Anime) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	inserted, err := config.Collection.InsertOne(ctx, anime)
	if err != nil {
		log.Println("Error inserting anime:", err)
		return err
	}
	fmt.Println("Inserted 1 anime in db with id:", inserted.InsertedID)
	return nil
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
