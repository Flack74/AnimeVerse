package router

import (
	"os"
	"time"

	controller "github.com/Flack74/mongoapi/controllers"
	middlewareAuth "github.com/Flack74/mongoapi/middleware"
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
	router.Get("/", controller.ServeFrontendHandler)
	router.Get("/api-home", controller.ServeHomeHandler)
	router.Get("/health", controller.HealthCheckHandler)
	
	// Public API routes (with optional auth for user-specific data)
	router.Route("/api", func(r chi.Router) {
		r.Use(middlewareAuth.OptionalSupabaseAuth)
		r.Get("/animes", controller.GetMyAllAnimesHandler)
		r.Get("/animes/filter", controller.FilterAnimesHandler)
		r.Get("/animes/trending", controller.GetTrendingAnimesHandler)
		r.Get("/animes/popular", controller.GetPopularAnimesHandler)
		r.Get("/anime/{animeName}", controller.GetAnimeByNameHandler)
	})
	
	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", controller.RegisterHandler)
		r.Post("/login", controller.LoginHandler)
		r.Post("/logout", controller.LogoutHandler)
		r.Get("/oauth", controller.SupabaseOAuthHandler)
	})
	
	// User routes (require Supabase auth)
	router.Route("/api/user", func(r chi.Router) {
		r.Use(middlewareAuth.SupabaseAuth)
		r.Get("/me", controller.GetCurrentUserHandler)
		r.Get("/stats", controller.GetUserStatsHandler)
		r.Post("/anime", controller.AddAnimeHandler)
		r.Put("/anime/{id}/status", controller.UpdateAnimeStatusHandler)
		r.Put("/anime/{id}/score", controller.UpdateAnimeScoreHandler)
		r.Delete("/anime/{id}", controller.RemoveAnimeHandler)
		r.Get("/search", controller.SearchAnimeHandler)
	})
	
	// Protected Admin API routes (Supabase auth + admin check)
	router.Route("/api/admin", func(r chi.Router) {
		r.Use(middlewareAuth.SupabaseAuth)
		r.Use(middlewareAuth.AdminOnly)
		r.Post("/anime", controller.CreateAnimeHandler)
		r.Post("/addmultipleanimes", controller.CreateMultipleAnimesHandler)
		r.Put("/anime/{id}", controller.UpdateAnimeHandler)
		r.Delete("/anime/{id}", controller.DeleteAnAnimeHandler)
		r.Delete("/deleteallanime", controller.DeleteEveryAnimesHandler)
		r.Post("/anime/{id}/episode/increment", controller.IncrementEpisodeHandler)
		r.Post("/anime/{id}/episode/decrement", controller.DecrementEpisodeHandler)
		r.Post("/anime/{id}/status/toggle", controller.ToggleStatusHandler)
		r.Post("/import/trending", controller.ImportTrendingHandler)
		r.Post("/import/seasonal", controller.ImportSeasonalHandler)
		r.Post("/backfill", controller.BackfillDataHandler)
	})
	
	// Legacy basic auth routes (fallback)
	legacyAuth := middlewareAuth.BasicAuth(
		os.Getenv("ADMIN_USERNAME"),
		os.Getenv("ADMIN_PASSWORD"),
	)
	if os.Getenv("ADMIN_USERNAME") == "" {
		legacyAuth = middlewareAuth.BasicAuth("admin", "admin123")
	}
	
	router.Route("/api/legacy", func(r chi.Router) {
		r.Use(legacyAuth)
		r.Post("/anime", controller.CreateAnimeHandler)
		r.Post("/addmultipleanimes", controller.CreateMultipleAnimesHandler)
		r.Put("/anime/{id}", controller.UpdateAnimeHandler)
		r.Delete("/anime/{id}", controller.DeleteAnAnimeHandler)
		r.Delete("/deleteallanime", controller.DeleteEveryAnimesHandler)
	})

	return router
}
