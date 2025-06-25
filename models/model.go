package model

import (
	"time"
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

// Season represents anime season
type Season string

const (
	Winter Season = "Winter"
	Spring Season = "Spring"
	Summer Season = "Summer"
	Fall   Season = "Fall"
)

// UserStats represents user statistics
type UserStats struct {
	TotalAnimes      int       `json:"total_animes" bson:"total_animes"`
	CompletedCount   int       `json:"completed_count" bson:"completed_count"`
	WatchingCount    int       `json:"watching_count" bson:"watching_count"`
	OnHoldCount      int       `json:"on_hold_count" bson:"on_hold_count"`
	DroppedCount     int       `json:"dropped_count" bson:"dropped_count"`
	PlanToWatchCount int       `json:"plan_to_watch_count" bson:"plan_to_watch_count"`
	LastUpdated      time.Time `json:"last_updated" bson:"last_updated"`
}

// User represents a user in the system
type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SupabaseID string             `json:"supabase_id" bson:"supabase_id"`
	Email      string             `json:"email" bson:"email"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Role       string             `json:"role" bson:"role"` // "user" or "admin"
	Stats      UserStats          `json:"stats" bson:"stats"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// Anime struct with improved validation and structure
type Anime struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id,omitempty" bson:"user_id,omitempty"` // Clerk user ID
	Name      string             `json:"name" bson:"name" validate:"required,min=1,max=200"`
	Type      AnimeType          `json:"type,omitempty" bson:"type,omitempty"`
	Score     int                `json:"score,omitempty" bson:"score,omitempty" validate:"min=0,max=10"` // Rating from 0 to 10
	Progress  Progress           `json:"progress,omitempty" bson:"progress,omitempty"`
	Status    WatchStatus        `json:"status,omitempty" bson:"status,omitempty"`
	Genre     []string           `json:"genre,omitempty" bson:"genre,omitempty"` // Genre tags
	Notes     string             `json:"notes,omitempty" bson:"notes,omitempty" validate:"max=500"` // Personal notes
	BannerUrl string             `json:"bannerUrl,omitempty" bson:"bannerUrl,omitempty"`
	ImageUrl  string             `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
	Year      int                `json:"year,omitempty" bson:"year,omitempty"`
	Season    Season             `json:"season,omitempty" bson:"season,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
