package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define AnimeType for type safety
type AnimeType string

const (
	SeriesType AnimeType = "TV"
	MovieType  AnimeType = "Movie"
	ONAType    AnimeType = "ONA"
)

// Define WatchStatus for user tracking
type WatchStatus string

const (
	Watching    WatchStatus = "watching"
	Completed   WatchStatus = "completed"
	OnHold      WatchStatus = "on-hold"
	Dropped     WatchStatus = "dropped"
	PlanToWatch WatchStatus = "plan-to-watch"
)

// Anime struct (Removed `Watched`)
type Anime struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Type     AnimeType          `json:"type,omitempty"`
	Score    int                `json:"score,omitempty"` // Rating from 1 to 10
	Progress struct {
		Watched int `json:"watched,omitempty"` // Episodes watched
		Total   int `json:"total,omitempty"`   // Total episodes (0 if unknown)
	} `json:"progress,omitempty"`
	Status WatchStatus `json:"status,omitempty"` // watching status
	Genre  []string    `json:"genre,omitempty"`  // Genre tags (e.g., "action", "isekai", "romance")
}
