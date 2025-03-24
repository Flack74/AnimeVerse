package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Flack74/mongoapi/config"
	"github.com/Flack74/mongoapi/router"
)

func main() {
	fmt.Println("Anime API")

	// Connect to MongoDB
	config.ConnectDB()

	r := router.Router()
	fmt.Println("Server is getting started...")
	fmt.Println("Listening at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
