package router

import (
	"net/http"
	"time"

	controller "animeverse/controllers"
	middlewareAuth "animeverse/middleware"

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
	router.Get("/old", controller.ServeOldFrontendHandler)
	router.Get("/health", controller.HealthCheckHandler)

	// Static files
	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public API routes (with optional auth for user-specific data)
	router.Route("/api", func(r chi.Router) {
		r.Use(middlewareAuth.OptionalSupabaseAuth)
		r.Get("/animes", controller.GetMyAllAnimesHandler)
		r.Get("/animes/filter", controller.FilterAnimesHandler)
		r.Get("/animes/trending", controller.GetTrendingAnimesHandler)
		r.Get("/animes/popular", controller.GetPopularAnimesHandler)
		r.Get("/animes/random", controller.GetRandomAnimeHandler)
		r.Get("/animes/top2025", controller.GetTop2025AnimesHandler)
		r.Get("/animes/preview", controller.GetPreviewAnimesHandler)
		r.Get("/animes/search", controller.SearchAnimesHandler)
		r.Get("/animes/spotlight", controller.GetSpotlightHandler)
		r.Get("/animes/top-rated-mixed", controller.GetTopRatedMixedHandler)
		r.Get("/animes/trending-fast", controller.GetTrendingFastHandler)
		r.Get("/anime/{animeName}", controller.GetAnimeByNameHandler)
		r.Get("/anime/fallback/{name}", controller.GetAnimeWithFallbackHandler)
		r.Get("/anime/themes", controller.GetAnimeThemesHandler)
		r.Get("/anime/hq-images", controller.GetHighQualityImagesHandler)
		r.Get("/anime/upgrade-images", controller.UpgradeImagesHandler)
		r.Get("/schedule/today", controller.GetScheduleHandler)

		// Fast loading endpoints with enhanced data
		r.Get("/fast/browse", controller.GetFastBrowseHandler)
		r.Get("/fast/top-rated", controller.GetFastTopRatedHandler)
		r.Get("/fast/search", controller.GetFastSearchHandler)

		// Backend-first endpoints (Database → Cache → External)
		r.Get("/backend/trending", controller.BackendFirstTrendingHandler)
		r.Get("/backend/browse", controller.BackendFirstBrowseHandler)
		r.Get("/simple/browse", controller.SimpleBrowseHandler)
		r.Get("/images/check", controller.CheckImagesHandler)
		r.Post("/images/save", controller.SaveImagesHandler)
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
		r.Post("/import/bulk", controller.BulkImportHandler)
		r.Post("/update/current", controller.UpdateCurrentSeasonHandler)
		r.Post("/backfill", controller.BackfillDataHandler)
		r.Post("/anime/{id}/enhance", controller.EnhanceAnime)
		r.Post("/anime/create-from-api", controller.CreateAnimeFromAPI)
		r.Get("/anime/{id}/enhanced", controller.GetEnhancedAnime)
	})

	router.Route("/api/legacy", func(r chi.Router) {
		r.Post("/anime", controller.CreateAnimeHandler)
		r.Post("/addmultipleanimes", controller.CreateMultipleAnimesHandler)
		r.Post("/import/bulk", controller.BulkImportHandler)
		r.Put("/anime/{id}", controller.UpdateAnimeHandler)
		r.Delete("/anime/{id}", controller.DeleteAnAnimeHandler)
		r.Delete("/deleteallanime", controller.DeleteEveryAnimesHandler)
	})

	return router
}
