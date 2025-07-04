package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JikanResponse represents Jikan API response
type JikanResponse struct {
	Data struct {
		Images struct {
			JPG struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"jpg"`
		} `json:"images"`
	} `json:"data"`
}

// AniListResponse represents AniList GraphQL response
type AniListResponse struct {
	Data struct {
		Media struct {
			BannerImage string `json:"bannerImage"`
			CoverImage  struct {
				Large string `json:"large"`
			} `json:"coverImage"`
		} `json:"Media"`
	} `json:"data"`
}

// ImageRequest represents a request to fetch images
type ImageRequest struct {
	MALID     int    `json:"mal_id,omitempty"`
	AniListID int    `json:"anilist_id,omitempty"`
	ImageUrl  string `json:"image_url,omitempty"`
	BannerUrl string `json:"banner_url,omitempty"`
}

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

// ImageCache represents cached image data
type ImageCache struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MALID       int                `json:"mal_id,omitempty" bson:"mal_id,omitempty"`
	AniListID   int                `json:"anilist_id,omitempty" bson:"anilist_id,omitempty"`
	ImageUrl    string             `json:"image_url,omitempty" bson:"image_url,omitempty"`
	BannerUrl   string             `json:"banner_url,omitempty" bson:"banner_url,omitempty"`
	LastUpdated time.Time          `json:"last_updated" bson:"last_updated"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

// Anime struct with detailed MAL-like information
type Anime struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name      string             `json:"name" bson:"name" validate:"required,min=1,max=200"`
	Type      AnimeType          `json:"type,omitempty" bson:"type,omitempty"`
	Score     float64            `json:"score,omitempty" bson:"score,omitempty"`
	Progress  Progress           `json:"progress,omitempty" bson:"progress,omitempty"`
	Status    WatchStatus        `json:"status,omitempty" bson:"status,omitempty"`
	Genre     []string           `json:"genre,omitempty" bson:"genre,omitempty"`
	Notes     string             `json:"notes,omitempty" bson:"notes,omitempty"`
	Synopsis  string             `json:"synopsis,omitempty" bson:"synopsis,omitempty"`
	BannerUrl string             `json:"bannerUrl,omitempty" bson:"bannerUrl,omitempty"`
	ImageUrl  string             `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
	MALID     int                `json:"mal_id,omitempty" bson:"mal_id,omitempty"`
	AniListID int                `json:"anilist_id,omitempty" bson:"anilist_id,omitempty"`
	Year      int                `json:"year,omitempty" bson:"year,omitempty"`
	Season    Season             `json:"season,omitempty" bson:"season,omitempty"`
	
	// Detailed Information
	AlternativeTitles AlternativeTitles `json:"alternative_titles,omitempty" bson:"alternative_titles,omitempty"`
	Information       AnimeInformation  `json:"information,omitempty" bson:"information,omitempty"`
	Statistics        AnimeStatistics   `json:"statistics,omitempty" bson:"statistics,omitempty"`
	Characters        []Character       `json:"characters,omitempty" bson:"characters,omitempty"`
	Staff             []StaffMember     `json:"staff,omitempty" bson:"staff,omitempty"`
	Themes            AnimeThemes       `json:"themes,omitempty" bson:"themes,omitempty"`
	Related           []RelatedAnime    `json:"related,omitempty" bson:"related,omitempty"`
	
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// AlternativeTitles represents alternative titles
type AlternativeTitles struct {
	Synonyms []string `json:"synonyms,omitempty" bson:"synonyms,omitempty"`
	Japanese string   `json:"japanese,omitempty" bson:"japanese,omitempty"`
	English  string   `json:"english,omitempty" bson:"english,omitempty"`
}

// AnimeInformation represents detailed anime information
type AnimeInformation struct {
	Episodes    int      `json:"episodes,omitempty" bson:"episodes,omitempty"`
	Status      string   `json:"status,omitempty" bson:"status,omitempty"`
	Aired       string   `json:"aired,omitempty" bson:"aired,omitempty"`
	Premiered   string   `json:"premiered,omitempty" bson:"premiered,omitempty"`
	Broadcast   string   `json:"broadcast,omitempty" bson:"broadcast,omitempty"`
	Producers   []string `json:"producers,omitempty" bson:"producers,omitempty"`
	Licensors   []string `json:"licensors,omitempty" bson:"licensors,omitempty"`
	Studios     []string `json:"studios,omitempty" bson:"studios,omitempty"`
	Source      string   `json:"source,omitempty" bson:"source,omitempty"`
	Duration    string   `json:"duration,omitempty" bson:"duration,omitempty"`
	Rating      string   `json:"rating,omitempty" bson:"rating,omitempty"`
}

// AnimeStatistics represents anime statistics
type AnimeStatistics struct {
	Score      float64 `json:"score,omitempty" bson:"score,omitempty"`
	ScoredBy   int     `json:"scored_by,omitempty" bson:"scored_by,omitempty"`
	Ranked     int     `json:"ranked,omitempty" bson:"ranked,omitempty"`
	Popularity int     `json:"popularity,omitempty" bson:"popularity,omitempty"`
	Members    int     `json:"members,omitempty" bson:"members,omitempty"`
	Favorites  int     `json:"favorites,omitempty" bson:"favorites,omitempty"`
}

// Character represents anime character
type Character struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Role        string `json:"role,omitempty" bson:"role,omitempty"`
	ImageUrl    string `json:"image_url,omitempty" bson:"image_url,omitempty"`
	VoiceActor  string `json:"voice_actor,omitempty" bson:"voice_actor,omitempty"`
	VAImageUrl  string `json:"va_image_url,omitempty" bson:"va_image_url,omitempty"`
}

// StaffMember represents anime staff
type StaffMember struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
	ImageUrl string `json:"image_url,omitempty" bson:"image_url,omitempty"`
}

// AnimeThemes represents opening and ending themes
type AnimeThemes struct {
	Openings []string `json:"openings,omitempty" bson:"openings,omitempty"`
	Endings  []string `json:"endings,omitempty" bson:"endings,omitempty"`
}

// RelatedAnime represents related anime
type RelatedAnime struct {
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	RelationType string `json:"relation_type,omitempty" bson:"relation_type,omitempty"`
	ImageUrl     string `json:"image_url,omitempty" bson:"image_url,omitempty"`
	MALID        int    `json:"mal_id,omitempty" bson:"mal_id,omitempty"`
}

// Episode represents an anime episode
type Episode struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AnimeID     primitive.ObjectID `json:"anime_id" bson:"anime_id"`
	Number      int                `json:"number" bson:"number"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Duration    int                `json:"duration,omitempty" bson:"duration,omitempty"` // in minutes
	AirDate     time.Time          `json:"air_date,omitempty" bson:"air_date,omitempty"`
	VideoUrl    string             `json:"video_url,omitempty" bson:"video_url,omitempty"`
	Thumbnail   string             `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

// Comment represents user comments on anime
type Comment struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AnimeID   primitive.ObjectID `json:"anime_id" bson:"anime_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	UserName  string             `json:"user_name" bson:"user_name"`
	Content   string             `json:"content" bson:"content" validate:"required,max=1000"`
	Spoiler   bool               `json:"spoiler,omitempty" bson:"spoiler,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
