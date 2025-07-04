package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "genuine-killdeer-55041.upstash.io:6379",
		Password: "AdcBAAIjcDEwYTBkNDE4ZjM3ZTk0ZGQ5OWU5ZDU5YTFlY2UwNmJkMXAxMA",
		DB:       0,
		TLSConfig: &tls.Config{},
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}
}

func Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, key, json, expiration).Err()
}

func Get(key string, dest interface{}) error {
	ctx := context.Background()
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func Delete(key string) error {
	ctx := context.Background()
	return RedisClient.Del(ctx, key).Err()
}

func Exists(key string) bool {
	ctx := context.Background()
	result, _ := RedisClient.Exists(ctx, key).Result()
	return result > 0
}