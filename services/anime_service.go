package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"animeverse/config"
	model "animeverse/models"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAnimeByName(name string) (*model.Anime, error) {
	collection := config.Collection
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

func InsertMultipleAnimes(animes []model.Anime) ([]interface{}, []string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var validAnimes []interface{}
	var duplicates []string

	for _, anime := range animes {
		if anime.Name == "" {
			continue
		}
		
		existing, err := FindAnimeByName(anime.Name)
		if err != nil {
			log.Printf("Error checking duplicate for %s: %v", anime.Name, err)
			continue
		}
		if existing != nil {
			duplicates = append(duplicates, anime.Name)
			continue
		}
		
		validAnimes = append(validAnimes, anime)
	}

	if len(validAnimes) == 0 {
		return nil, duplicates, fmt.Errorf("no valid animes to insert")
	}

	result, err := config.Collection.InsertMany(ctx, validAnimes)
	if err != nil {
		log.Println("Error inserting multiple animes:", err)
		return nil, duplicates, err
	}

	fmt.Printf("Inserted %d animes in db\n", len(result.InsertedIDs))
	return result.InsertedIDs, duplicates, nil
}

func UpdateAnime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	delete(updates, "_id")

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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"matched":  result.MatchedCount,
		"modified": result.ModifiedCount,
		"id":       idStr,
	})
}

func DeleteOneAnime(animeId string) bool {
	id, err := primitive.ObjectIDFromHex(animeId)
	if err != nil {
		log.Println("Invalid ObjectID:", err)
		return false
	}

	filter := bson.M{"_id": id}
	result, err := config.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Error deleting anime:", err)
		return false
	}
	
	if result.DeletedCount == 0 {
		log.Println("No anime found with ID:", animeId)
		return false
	}
	
	fmt.Println("Anime deleted with count:", result.DeletedCount)
	return true
}

func DeleteAllAnime() int64 {
	deleteResult, err := config.Collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Println("Error deleting all animes:", err)
		return 0
	}
	fmt.Println("Number of animes deleted:", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

func GetAllAnimes() []primitive.M {
	cur, err := config.Collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Println("Error fetching animes:", err)
		return []primitive.M{} // Return empty slice instead of nil
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
	log.Printf("GetAllAnimes: Found %d total animes", len(animes))
	return animes
}

func FilterAnimes(search, genre, year, season, format, status, userID string) []primitive.M {
	filter := bson.M{}
	
	// Add user filter if provided (for user-specific data)
	if userID != "" {
		filter["user_id"] = userID
	}
	
	// Build filter with proper field matching
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}
	if genre != "" {
		// Handle both string and array genre fields
		filter["$or"] = []bson.M{
			{"genre": bson.M{"$regex": genre, "$options": "i"}}, // String genre
			{"genre": bson.M{"$in": []string{genre}}},              // Array genre
		}
	}
	if year != "" {
		if yearInt, err := strconv.Atoi(year); err == nil {
			filter["year"] = yearInt
		}
	}
	if season != "" {
		filter["season"] = bson.M{"$regex": "^" + season + "$", "$options": "i"}
	}
	if format != "" {
		filter["type"] = bson.M{"$regex": "^" + format + "$", "$options": "i"}
	}
	if status != "" {
		filter["status"] = bson.M{"$regex": "^" + status + "$", "$options": "i"}
	}

	// Log filter for debugging
	log.Printf("Filter query: %+v", filter)
	
	cur, err := config.Collection.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error filtering animes:", err)
		return []primitive.M{} // Return empty slice instead of nil
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
	
	log.Printf("Found %d animes with filter", len(animes))
	return animes
}

func GetTrendingAnimes() []primitive.M {
	filter := bson.M{}
	opts := options.Find().SetSort(bson.D{{"score", -1}}).SetLimit(5)
	
	cur, err := config.Collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println("Error fetching trending animes:", err)
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

func GetPopularAnimes() []primitive.M {
	filter := bson.M{"status": "completed"}
	opts := options.Find().SetSort(bson.D{{"score", -1}}).SetLimit(5)
	
	cur, err := config.Collection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println("Error fetching popular animes:", err)
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

// GetRandomAnime returns a random anime
func GetRandomAnime() *model.Anime {
	count, err := config.Collection.CountDocuments(context.Background(), bson.M{})
	if err != nil || count == 0 {
		return nil
	}
	
	randomIndex := rand.Int63n(count)
	opts := options.FindOne().SetSkip(randomIndex)
	
	var anime model.Anime
	err = config.Collection.FindOne(context.Background(), bson.M{}, opts).Decode(&anime)
	if err != nil {
		return nil
	}
	
	return &anime
}

// GetDailySchedule returns scheduled anime for a given day
func GetDailySchedule(day string) []map[string]interface{} {
	// Mock schedule data - in real implementation, this would come from a schedule collection
	schedules := map[string][]map[string]interface{}{
		"Monday": {
			{"name": "Attack on Titan", "time": "15:45", "episode": "Episode 12"},
			{"name": "Demon Slayer", "time": "16:30", "episode": "Episode 8"},
			{"name": "One Piece", "time": "17:00", "episode": "Episode 1045"},
		},
		"Tuesday": {
			{"name": "Jujutsu Kaisen", "time": "15:30", "episode": "Episode 15"},
			{"name": "My Hero Academia", "time": "16:00", "episode": "Episode 22"},
		},
		"Wednesday": {
			{"name": "Chainsaw Man", "time": "15:45", "episode": "Episode 9"},
			{"name": "Spy x Family", "time": "16:15", "episode": "Episode 11"},
		},
		"Thursday": {
			{"name": "Tokyo Revengers", "time": "15:30", "episode": "Episode 18"},
			{"name": "Blue Lock", "time": "16:45", "episode": "Episode 14"},
		},
		"Friday": {
			{"name": "Mob Psycho 100", "time": "15:00", "episode": "Episode 7"},
			{"name": "One Punch Man", "time": "16:30", "episode": "Episode 5"},
		},
		"Saturday": {
			{"name": "Dragon Ball Super", "time": "14:00", "episode": "Episode 131"},
			{"name": "Naruto Shippuden", "time": "15:30", "episode": "Episode 720"},
		},
		"Sunday": {
			{"name": "Bleach", "time": "14:30", "episode": "Episode 366"},
			{"name": "Hunter x Hunter", "time": "16:00", "episode": "Episode 148"},
		},
	}
	
	if schedule, exists := schedules[day]; exists {
		return schedule
	}
	
	return []map[string]interface{}{}
}

// GetTop2025Animes returns top rated anime from 2024-2025
func GetTop2025Animes() []primitive.M {
	// Multi-stage approach for best results
	filters := []bson.M{
		// Stage 1: 2024-2025 high rated
		{"$and": []bson.M{{"year": bson.M{"$gte": 2024}}, {"score": bson.M{"$gte": 8}}}},
		// Stage 2: 2023-2024 very high rated
		{"$and": []bson.M{{"year": bson.M{"$gte": 2023}}, {"score": bson.M{"$gte": 9}}}},
		// Stage 3: Any high rated with images
		{"$and": []bson.M{{"score": bson.M{"$gte": 8}}, {"imageUrl": bson.M{"$ne": ""}}}},
		// Stage 4: Popular genres
		{"$and": []bson.M{{"score": bson.M{"$gte": 7}}, {"genre": bson.M{"$in": []string{"Action", "Adventure", "Fantasy", "Shounen"}}}}},
	}
	
	var animes []primitive.M
	seenNames := make(map[string]bool)
	
	for _, filter := range filters {
		if len(animes) >= 5 {
			break
		}
		
		opts := options.Find().SetSort(bson.D{{"score", -1}, {"year", -1}}).SetLimit(10)
		cur, err := config.Collection.Find(context.Background(), filter, opts)
		if err != nil {
			continue
		}
		
		for cur.Next(context.Background()) && len(animes) < 5 {
			var anime bson.M
			if err := cur.Decode(&anime); err != nil {
				continue
			}
			
			name, ok := anime["name"].(string)
			if !ok || seenNames[name] {
				continue
			}
			
			seenNames[name] = true
			animes = append(animes, anime)
		}
		cur.Close(context.Background())
	}
	
	log.Printf("GetTop2025Animes: Found %d anime", len(animes))
	return animes
}

// GetPreviewAnimes returns a preview of anime for the main page
func GetPreviewAnimes() []primitive.M {
	// Get diverse mix of anime
	pipeline := []bson.M{
		{"$match": bson.M{"score": bson.M{"$gte": 6}}},
		{"$sort": bson.M{"score": -1, "year": -1}},
		{"$limit": 12},
	}
	
	cur, err := config.Collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Println("Error fetching preview animes:", err)
		// Fallback to simple query
		return GetAllAnimes()[:12]
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
	
	log.Printf("GetPreviewAnimes: Found %d anime", len(animes))
	return animes
}

// SearchAnimes searches for anime by name with fuzzy matching
func SearchAnimes(query string) []primitive.M {
	// Multi-stage search for better results
	filters := []bson.M{
		// Exact match
		{"name": bson.M{"$regex": "^" + query + "$", "$options": "i"}},
		// Starts with
		{"name": bson.M{"$regex": "^" + query, "$options": "i"}},
		// Contains
		{"name": bson.M{"$regex": query, "$options": "i"}},
		// Synonyms/notes search
		{"notes": bson.M{"$regex": query, "$options": "i"}},
	}
	
	var allResults []primitive.M
	seenNames := make(map[string]bool)
	
	for _, filter := range filters {
		if len(allResults) >= 20 { // Limit total results
			break
		}
		
		opts := options.Find().SetSort(bson.D{{"score", -1}}).SetLimit(10)
		cur, err := config.Collection.Find(context.Background(), filter, opts)
		if err != nil {
			continue
		}
		
		for cur.Next(context.Background()) && len(allResults) < 20 {
			var anime bson.M
			if err := cur.Decode(&anime); err != nil {
				continue
			}
			
			name, ok := anime["name"].(string)
			if !ok || seenNames[name] {
				continue
			}
			
			seenNames[name] = true
			allResults = append(allResults, anime)
		}
		cur.Close(context.Background())
	}
	
	log.Printf("SearchAnimes: Found %d results for query '%s'", len(allResults), query)
	return allResults
}
