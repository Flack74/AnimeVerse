package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"animeverse/cache"
	"animeverse/config"
	"animeverse/router"
	"animeverse/services"
)

func main() {
	fmt.Println("🎌 AnimeVerse API - Starting...")

	// Connect to MongoDB
	config.ConnectDB()
	
	// Initialize Redis cache
	cache.InitRedis()
	
	// Initialize image collection
	services.InitImageCollection()

	// Setup router
	r := router.Router()

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("🚀 Server starting on port %s\n", port)
		fmt.Printf("🌐 Access the API at: http://localhost:%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server exited gracefully")
}
