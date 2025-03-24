package router

import (
	controller "github.com/Flack74/mongoapi/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", controller.ServeHome)
	router.HandleFunc("/api/animes", controller.GetMyAllAnimes).Methods("GET")
	router.HandleFunc("/api/anime", controller.CreateAnime).Methods("POST")
	router.HandleFunc("/api/anime/{id}", controller.UpdateAnime).Methods("PUT")
	router.HandleFunc("/api/anime/{id}", controller.DeleteAnAnime).Methods("DELETE")
	router.HandleFunc("/api/deleteallanime", controller.DeleteEveryAnimes).Methods("DELETE")

	return router
}
