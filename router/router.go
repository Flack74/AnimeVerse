package router

import (
	"time"

	controller "github.com/Flack74/mongoapi/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Router() *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Compress(5))

	// CORS configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	router.Get("/", controller.ServeHomeHandler)
	router.Get("/health", controller.HealthCheckHandler)
	
	// API routes
	router.Route("/api", func(r chi.Router) {
		r.Get("/animes", controller.GetMyAllAnimesHandler)
		r.Get("/anime/{animeName}", controller.GetAnimeByNameHandler)
		r.Post("/anime", controller.CreateAnimeHandler)
		r.Put("/anime/{id}", controller.UpdateAnimeHandler)
		r.Delete("/anime/{id}", controller.DeleteAnAnimeHandler)
		r.Delete("/deleteallanime", controller.DeleteEveryAnimesHandler)
	})

	return router
}
