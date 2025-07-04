package services

import (
	"context"
	"time"

	"animeverse/cache"
)

const (
	RATE_LIMIT_PREFIX = "rate_limit:"
	DEFAULT_WINDOW    = 60 * time.Second
	DEFAULT_LIMIT     = 100
)

type RateLimiter struct {
	Window time.Duration
	Limit  int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		Window: window,
		Limit:  limit,
	}
}

func (rl *RateLimiter) Allow(key string) (bool, error) {
	now := time.Now().Unix()
	windowStart := now - int64(rl.Window.Seconds())
	
	// Use Lua script for atomic operations
	luaScript := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local now = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		
		-- Remove old entries
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start)
		
		-- Count current requests
		local current = redis.call('ZCARD', key)
		
		if current < limit then
			-- Add current request
			redis.call('ZADD', key, now, now)
			redis.call('EXPIRE', key, 60)
			return 1
		else
			return 0
		end
	`
	
	ctx := context.Background()
	redisKey := RATE_LIMIT_PREFIX + key
	result := cache.RedisClient.Eval(ctx, luaScript, []string{redisKey}, windowStart, now, rl.Limit)
	
	allowed, err := result.Int()
	if err != nil {
		return false, err
	}
	
	return allowed == 1, nil
}

func (rl *RateLimiter) GetRemaining(key string) (int, error) {
	ctx := context.Background()
	redisKey := RATE_LIMIT_PREFIX + key
	current, err := cache.RedisClient.ZCard(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	
	remaining := rl.Limit - int(current)
	if remaining < 0 {
		remaining = 0
	}
	
	return remaining, nil
}

// Predefined rate limiters
var (
	LoginLimiter    = NewRateLimiter(15*time.Minute, 5)   // 5 attempts per 15 minutes
	SearchLimiter   = NewRateLimiter(1*time.Minute, 30)   // 30 searches per minute
	APILimiter      = NewRateLimiter(1*time.Minute, 100)  // 100 API calls per minute
	CommentLimiter  = NewRateLimiter(5*time.Minute, 10)   // 10 comments per 5 minutes
)