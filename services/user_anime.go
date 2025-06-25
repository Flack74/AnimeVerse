package services

import (
	"context"
	"time"

	"github.com/Flack74/mongoapi/config"
	model "github.com/Flack74/mongoapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAnimeToUserList(userID, animeName string, status model.WatchStatus) (*model.Anime, error) {
	// First try to find anime in database
	existingAnime, err := FindAnimeByName(animeName)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	var anime model.Anime

	if existingAnime != nil {
		// Create user's copy of existing anime
		anime = *existingAnime
		anime.ID = primitive.NewObjectID()
		anime.UserID = userID
		anime.Status = status
		anime.Score = 0
		anime.Progress = model.Progress{Watched: 0, Total: anime.Progress.Total}
		anime.CreatedAt = time.Now()
		anime.UpdatedAt = time.Now()
	} else {
		// Try to fetch from external APIs
		if err := ImportFromAniList(animeName); err == nil {
			// Try to find it again after import
			if importedAnime, err := FindAnimeByName(animeName); err == nil {
				anime = *importedAnime
				anime.ID = primitive.NewObjectID()
				anime.UserID = userID
				anime.Status = status
				anime.Score = 0
				anime.Progress = model.Progress{Watched: 0, Total: anime.Progress.Total}
				anime.CreatedAt = time.Now()
				anime.UpdatedAt = time.Now()
			}
		}

		// If still not found, create basic entry
		if anime.Name == "" {
			anime = model.Anime{
				UserID:    userID,
				Name:      animeName,
				Type:      "TV",
				Score:     0,
				Progress:  model.Progress{Watched: 0, Total: 0},
				Status:    status,
				Genre:     []string{},
				Notes:     "Added by user",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
	}

	// Insert user's anime
	err = InsertOneAnime(anime)
	if err != nil {
		return nil, err
	}

	return &anime, nil
}

func UpdateAnimeStatus(userID, animeID string, status model.WatchStatus) (*model.Anime, error) {
	objID, err := primitive.ObjectIDFromHex(animeID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":     objID,
		"user_id": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err = config.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Return updated anime
	var anime model.Anime
	err = config.Collection.FindOne(context.Background(), filter).Decode(&anime)
	return &anime, err
}

func UpdateAnimeScore(userID, animeID string, score int) (*model.Anime, error) {
	objID, err := primitive.ObjectIDFromHex(animeID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":     objID,
		"user_id": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"score":      score,
			"updated_at": time.Now(),
		},
	}

	_, err = config.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	var anime model.Anime
	err = config.Collection.FindOne(context.Background(), filter).Decode(&anime)
	return &anime, err
}

func RemoveAnimeFromUserList(userID, animeID string) error {
	objID, err := primitive.ObjectIDFromHex(animeID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":     objID,
		"user_id": userID,
	}

	_, err = config.Collection.DeleteOne(context.Background(), filter)
	return err
}

func SearchAndAddAnime(userID, query string) ([]map[string]interface{}, error) {
	// Search external APIs for anime
	results := []map[string]interface{}{}

	// Try Jikan API
	if jikanResults, err := SearchJikanAPI(query); err == nil {
		for _, anime := range jikanResults {
			results = append(results, map[string]interface{}{
				"name":     anime.Title,
				"year":     anime.Year,
				"type":     anime.Type,
				"genres":   anime.Genres,
				"image":    anime.Images.JPG.ImageURL,
				"synopsis": anime.Synopsis,
				"source":   "jikan",
			})
		}
	}

	return results, nil
}