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

func CreateOrUpdateUser(supabaseID, email, name string) (*model.User, error) {
	ctx := context.Background()
	
	// Check if user exists
	var existingUser model.User
	err := config.UserCollection.FindOne(ctx, bson.M{"supabase_id": supabaseID}).Decode(&existingUser)
	
	if err == mongo.ErrNoDocuments {
		// Create new user with initial stats
		user := model.User{
			SupabaseID: supabaseID,
			Email:     email,
			Name:      name,
			Role:      "user", // Default role
			Stats: model.UserStats{
				TotalAnimes:      0,
				CompletedCount:   0,
				WatchingCount:    0,
				OnHoldCount:      0,
				DroppedCount:     0,
				PlanToWatchCount: 0,
				LastUpdated:      time.Now(),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		result, err := config.UserCollection.InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}
		
		user.ID = result.InsertedID.(primitive.ObjectID)
		return &user, nil
	} else if err != nil {
		return nil, err
	}
	
	// Update existing user
	existingUser.Email = email
	existingUser.Name = name
	existingUser.UpdatedAt = time.Now()
	
	_, err = config.UserCollection.UpdateOne(
		ctx,
		bson.M{"supabase_id": supabaseID},
		bson.M{"$set": bson.M{
			"email":      email,
			"name":       name,
			"updated_at": time.Now(),
		}},
	)
	
	return &existingUser, err
}

func GetUserBySupabaseID(supabaseID string) (*model.User, error) {
	var user model.User
	err := config.UserCollection.FindOne(context.Background(), bson.M{"supabase_id": supabaseID}).Decode(&user)
	return &user, err
}

func SetUserRole(supabaseID, role string) error {
	_, err := config.UserCollection.UpdateOne(
		context.Background(),
		bson.M{"supabase_id": supabaseID},
		bson.M{"$set": bson.M{
			"role":       role,
			"updated_at": time.Now(),
		}},
	)
	return err
}

func IsUserAdmin(supabaseID string) bool {
	user, err := GetUserBySupabaseID(supabaseID)
	if err != nil {
		return false
	}
	return user.Role == "admin"
}

func GetUserStats(supabaseID string) (*model.UserStats, error) {
	user, err := GetUserBySupabaseID(supabaseID)
	if err != nil {
		return nil, err
	}
	return &user.Stats, nil
}

func UpdateUserStats(supabaseID string) error {
	ctx := context.Background()
	
	// Count anime by status for this user
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": supabaseID}},
		{"$group": bson.M{
			"_id": "$status",
			"count": bson.M{"$sum": 1},
		}},
	}
	
	cur, err := config.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)
	
	stats := model.UserStats{
		LastUpdated: time.Now(),
	}
	
	for cur.Next(ctx) {
		var result bson.M
		if err := cur.Decode(&result); err != nil {
			continue
		}
		
		status := result["_id"].(string)
		count := int(result["count"].(int32))
		
		switch status {
		case "completed":
			stats.CompletedCount = count
		case "watching":
			stats.WatchingCount = count
		case "on-hold":
			stats.OnHoldCount = count
		case "dropped":
			stats.DroppedCount = count
		case "plan-to-watch":
			stats.PlanToWatchCount = count
		}
		stats.TotalAnimes += count
	}
	
	// Update user stats in MongoDB
	_, err = config.UserCollection.UpdateOne(
		ctx,
		bson.M{"supabase_id": supabaseID},
		bson.M{"$set": bson.M{"stats": stats}},
	)
	
	return err
}