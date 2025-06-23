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

// Progress represents anime watching progress
type Progress struct {
	Watched int `json:"watched,omitempty" bson:"watched,omitempty"` // Episodes watched
	Total   int `json:"total,omitempty" bson:"total,omitempty"`     // Total episodes (0 if unknown)
}

// Anime struct with improved validation and structure
type Anime struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" validate:"required,min=1,max=200"`
	Type      AnimeType          `json:"type,omitempty" bson:"type,omitempty"`
	Score     int                `json:"score,omitempty" bson:"score,omitempty" validate:"min=0,max=10"` // Rating from 0 to 10
	Progress  Progress           `json:"progress,omitempty" bson:"progress,omitempty"`
	Status    WatchStatus        `json:"status,omitempty" bson:"status,omitempty"`
	Genre     []string           `json:"genre,omitempty" bson:"genre,omitempty"` // Genre tags
	Notes     string             `json:"notes,omitempty" bson:"notes,omitempty" validate:"max=500"` // Personal notes
	BannerUrl string             `json:"bannerUrl,omitempty" bson:"bannerUrl,omitempty"`
	ImageUrl  string             `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
}
