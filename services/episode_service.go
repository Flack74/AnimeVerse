package services

import (
	"context"

	"github.com/Flack74/mongoapi/config"
	model "github.com/Flack74/mongoapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IncrementEpisode(animeId string) (*model.Anime, error) {
	id, err := primitive.ObjectIDFromHex(animeId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"progress.watched": 1}}

	_, err = config.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	var anime model.Anime
	err = config.Collection.FindOne(context.Background(), filter).Decode(&anime)
	return &anime, err
}

func DecrementEpisode(animeId string) (*model.Anime, error) {
	id, err := primitive.ObjectIDFromHex(animeId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"progress.watched": -1}}

	_, err = config.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	var anime model.Anime
	err = config.Collection.FindOne(context.Background(), filter).Decode(&anime)
	return &anime, err
}

func ToggleStatus(animeId string) (*model.Anime, error) {
	id, err := primitive.ObjectIDFromHex(animeId)
	if err != nil {
		return nil, err
	}

	var anime model.Anime
	filter := bson.M{"_id": id}
	err = config.Collection.FindOne(context.Background(), filter).Decode(&anime)
	if err != nil {
		return nil, err
	}

	newStatus := "watching"
	switch anime.Status {
	case "watching":
		newStatus = "completed"
	case "completed":
		newStatus = "on-hold"
	case "on-hold":
		newStatus = "dropped"
	case "dropped":
		newStatus = "plan-to-watch"
	default:
		newStatus = "watching"
	}

	update := bson.M{"$set": bson.M{"status": newStatus}}
	_, err = config.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	anime.Status = model.WatchStatus(newStatus)
	return &anime, nil
}