package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"animeverse/cache"
)

const (
	JOB_QUEUE_KEY = "job_queue"
	JOB_PROCESSING_KEY = "job_processing"
)

type JobType string

const (
	FetchAnimeImages JobType = "fetch_anime_images"
	UpdateAnimeData  JobType = "update_anime_data"
	SendNotification JobType = "send_notification"
)

type Job struct {
	ID        string                 `json:"id"`
	Type      JobType                `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
	Attempts  int                    `json:"attempts"`
	MaxRetries int                   `json:"max_retries"`
}

type JobQueue struct {
	ctx context.Context
}

func NewJobQueue() *JobQueue {
	return &JobQueue{
		ctx: context.Background(),
	}
}

func (jq *JobQueue) Enqueue(jobType JobType, payload map[string]interface{}) error {
	job := Job{
		ID:         fmt.Sprintf("%s_%d", jobType, time.Now().UnixNano()),
		Type:       jobType,
		Payload:    payload,
		CreatedAt:  time.Now(),
		Attempts:   0,
		MaxRetries: 3,
	}
	
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}
	
	return cache.RedisClient.LPush(jq.ctx, JOB_QUEUE_KEY, jobData).Err()
}

func (jq *JobQueue) Dequeue() (*Job, error) {
	result, err := cache.RedisClient.BRPop(jq.ctx, 5*time.Second, JOB_QUEUE_KEY).Result()
	if err != nil {
		return nil, err
	}
	
	if len(result) < 2 {
		return nil, fmt.Errorf("invalid job data")
	}
	
	var job Job
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, err
	}
	
	// Mark job as processing
	processingData, _ := json.Marshal(job)
	cache.RedisClient.HSet(jq.ctx, JOB_PROCESSING_KEY, job.ID, processingData)
	
	return &job, nil
}

func (jq *JobQueue) CompleteJob(jobID string) error {
	return cache.RedisClient.HDel(jq.ctx, JOB_PROCESSING_KEY, jobID).Err()
}

func (jq *JobQueue) RetryJob(job *Job) error {
	job.Attempts++
	
	if job.Attempts >= job.MaxRetries {
		log.Printf("Job %s exceeded max retries, discarding", job.ID)
		return jq.CompleteJob(job.ID)
	}
	
	// Re-queue with delay
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}
	
	// Add delay based on attempt count
	delay := time.Duration(job.Attempts) * 30 * time.Second
	
	go func() {
		time.Sleep(delay)
		cache.RedisClient.LPush(jq.ctx, JOB_QUEUE_KEY, jobData)
		jq.CompleteJob(job.ID)
	}()
	
	return nil
}

func (jq *JobQueue) ProcessJobs() {
	log.Println("Starting job queue processor...")
	
	for {
		job, err := jq.Dequeue()
		if err != nil {
			continue
		}
		
		log.Printf("Processing job: %s (type: %s)", job.ID, job.Type)
		
		if err := jq.processJob(job); err != nil {
			log.Printf("Job %s failed: %v", job.ID, err)
			jq.RetryJob(job)
		} else {
			log.Printf("Job %s completed successfully", job.ID)
			jq.CompleteJob(job.ID)
		}
	}
}

func (jq *JobQueue) processJob(job *Job) error {
	switch job.Type {
	case FetchAnimeImages:
		return jq.processFetchAnimeImages(job)
	case UpdateAnimeData:
		return jq.processUpdateAnimeData(job)
	case SendNotification:
		return jq.processSendNotification(job)
	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}

func (jq *JobQueue) processFetchAnimeImages(job *Job) error {
	animeName, ok := job.Payload["anime_name"].(string)
	if !ok {
		return fmt.Errorf("invalid anime_name in payload")
	}
	
	animeID, ok := job.Payload["anime_id"].(string)
	if !ok {
		return fmt.Errorf("invalid anime_id in payload")
	}
	
	// Update anime with new images
	return UpdateAnimeWithAPIData(animeID, animeName)
}

func (jq *JobQueue) processUpdateAnimeData(job *Job) error {
	animeName, ok := job.Payload["anime_name"].(string)
	if !ok {
		return fmt.Errorf("invalid anime_name in payload")
	}
	
	// Enhance anime data from API
	_, err := EnhanceAnimeFromAPI(animeName)
	return err
}

func (jq *JobQueue) processSendNotification(job *Job) error {
	message, ok := job.Payload["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in payload")
	}
	
	log.Printf("Notification: %s", message)
	return nil
}

func (jq *JobQueue) fetchHighQualityImages(animeName string) (string, string, error) {
	// Try multiple image sources for better quality
	sources := []func(string) (string, string, error){
		jq.fetchFromAniList,
		jq.fetchFromJikan,
		jq.fetchFromTMDB,
	}
	
	for _, source := range sources {
		imageUrl, bannerUrl, err := source(animeName)
		if err == nil && imageUrl != "" {
			return imageUrl, bannerUrl, nil
		}
	}
	
	return "", "", fmt.Errorf("no high-quality images found")
}

func (jq *JobQueue) fetchFromAniList(animeName string) (string, string, error) {
	// Implementation for AniList image fetching
	return "", "", fmt.Errorf("not implemented")
}

func (jq *JobQueue) fetchFromJikan(animeName string) (string, string, error) {
	// Implementation for Jikan image fetching
	return "", "", fmt.Errorf("not implemented")
}

func (jq *JobQueue) fetchFromTMDB(animeName string) (string, string, error) {
	// Implementation for TMDB image fetching
	return "", "", fmt.Errorf("not implemented")
}

// Global job queue instance
var GlobalJobQueue = NewJobQueue()

// Helper functions for common job types
func EnqueueImageFetch(animeID, animeName string) error {
	return GlobalJobQueue.Enqueue(FetchAnimeImages, map[string]interface{}{
		"anime_id":   animeID,
		"anime_name": animeName,
	})
}

func EnqueueDataUpdate(animeName string) error {
	return GlobalJobQueue.Enqueue(UpdateAnimeData, map[string]interface{}{
		"anime_name": animeName,
	})
}

func EnqueueNotification(message string) error {
	return GlobalJobQueue.Enqueue(SendNotification, map[string]interface{}{
		"message": message,
	})
}