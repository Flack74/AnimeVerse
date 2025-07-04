package services

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"animeverse/cache"
	"animeverse/config"
	"animeverse/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const SEARCH_CACHE_PREFIX = "search:"

type SearchResult struct {
	Anime *models.Anime `json:"anime"`
	Score float64       `json:"score"`
}

type SearchService struct {
	CacheDuration time.Duration
}

func NewSearchService() *SearchService {
	return &SearchService{
		CacheDuration: 10 * time.Minute,
	}
}

func (s *SearchService) SearchAnime(query string, limit int) ([]*models.Anime, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("%s%s:%d", SEARCH_CACHE_PREFIX, strings.ToLower(query), limit)
	
	var cachedResults []*models.Anime
	if err := cache.Get(cacheKey, &cachedResults); err == nil {
		return cachedResults, nil
	}
	
	// Perform search
	results, err := s.performSearch(query, limit)
	if err != nil {
		return nil, err
	}
	
	// Cache results
	cache.Set(cacheKey, results, s.CacheDuration)
	
	return results, nil
}

func (s *SearchService) performSearch(query string, limit int) ([]*models.Anime, error) {
	ctx := context.Background()
	
	// Multi-stage search strategy
	var allResults []SearchResult
	
	// 1. Exact name match (highest priority)
	exactResults, _ := s.searchExact(ctx, query)
	for _, anime := range exactResults {
		allResults = append(allResults, SearchResult{Anime: anime, Score: 100.0})
	}
	
	// 2. Fuzzy name match
	fuzzyResults, _ := s.searchFuzzy(ctx, query)
	for _, anime := range fuzzyResults {
		score := s.calculateSimilarity(query, anime.Name)
		allResults = append(allResults, SearchResult{Anime: anime, Score: score * 80})
	}
	
	// 3. Genre match
	genreResults, _ := s.searchByGenre(ctx, query)
	for _, anime := range genreResults {
		allResults = append(allResults, SearchResult{Anime: anime, Score: 60.0})
	}
	
	// Remove duplicates and sort by score
	uniqueResults := s.removeDuplicates(allResults)
	sort.Slice(uniqueResults, func(i, j int) bool {
		if uniqueResults[i].Score == uniqueResults[j].Score {
			// Secondary sort by popularity (score)
			return uniqueResults[i].Anime.Score > uniqueResults[j].Anime.Score
		}
		return uniqueResults[i].Score > uniqueResults[j].Score
	})
	
	// Extract anime and limit results
	var results []*models.Anime
	for i, result := range uniqueResults {
		if i >= limit {
			break
		}
		results = append(results, result.Anime)
	}
	
	return results, nil
}

func (s *SearchService) searchExact(ctx context.Context, query string) ([]*models.Anime, error) {
	filter := bson.M{
		"name": bson.M{"$regex": "^" + query + "$", "$options": "i"},
	}
	
	return s.executeSearch(ctx, filter, 5)
}

func (s *SearchService) searchFuzzy(ctx context.Context, query string) ([]*models.Anime, error) {
	filter := bson.M{
		"name": bson.M{"$regex": query, "$options": "i"},
	}
	
	return s.executeSearch(ctx, filter, 20)
}

func (s *SearchService) searchByGenre(ctx context.Context, query string) ([]*models.Anime, error) {
	filter := bson.M{
		"genre": bson.M{"$in": []string{query}},
	}
	
	return s.executeSearch(ctx, filter, 10)
}

func (s *SearchService) executeSearch(ctx context.Context, filter bson.M, limit int) ([]*models.Anime, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{"score", -1}})
	
	cursor, err := config.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var animes []*models.Anime
	for cursor.Next(ctx) {
		var anime models.Anime
		if err := cursor.Decode(&anime); err != nil {
			continue
		}
		animes = append(animes, &anime)
	}
	
	return animes, nil
}

func (s *SearchService) calculateSimilarity(query, target string) float64 {
	query = strings.ToLower(query)
	target = strings.ToLower(target)
	
	// Simple similarity calculation
	if query == target {
		return 1.0
	}
	
	if strings.Contains(target, query) {
		return 0.8
	}
	
	// Levenshtein distance approximation
	return s.levenshteinSimilarity(query, target)
}

func (s *SearchService) levenshteinSimilarity(s1, s2 string) float64 {
	if len(s1) == 0 {
		return float64(len(s2))
	}
	if len(s2) == 0 {
		return float64(len(s1))
	}
	
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}
	
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			
			deletion := matrix[i-1][j] + 1
			insertion := matrix[i][j-1] + 1
			substitution := matrix[i-1][j-1] + cost
			
			matrix[i][j] = min(deletion, min(insertion, substitution))
		}
	}
	
	distance := matrix[len(s1)][len(s2)]
	maxLen := max(len(s1), len(s2))
	
	return 1.0 - float64(distance)/float64(maxLen)
}

func (s *SearchService) removeDuplicates(results []SearchResult) []SearchResult {
	seen := make(map[string]bool)
	var unique []SearchResult
	
	for _, result := range results {
		if !seen[result.Anime.ID.Hex()] {
			seen[result.Anime.ID.Hex()] = true
			unique = append(unique, result)
		}
	}
	
	return unique
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}