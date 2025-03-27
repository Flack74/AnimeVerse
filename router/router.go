package router

import (
	controller "github.com/Flack74/mongoapi/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", controller.ServeHomeHandler)
	router.HandleFunc("/api/animes", controller.GetMyAllAnimesHandler).Methods("GET")
	router.HandleFunc("/api/anime/{animeName}", controller.GetAnimeByNameHandler)
	router.HandleFunc("/api/anime", controller.CreateAnimeHandler).Methods("POST")
	router.HandleFunc("/api/anime/{id}", controller.UpdateAnimeHandler).Methods("PUT")
	router.HandleFunc("/api/anime/{id}", controller.DeleteAnAnimeHandler).Methods("DELETE")
	router.HandleFunc("/api/deleteallanime", controller.DeleteEveryAnimesHandler).Methods("DELETE")

	return router
}
