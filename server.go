package main

import (
	"log"
	"net/http"

	gameHandlers "github.com/LYSingD/go-games-rest-api/gameHandlers"
)

func main() {
	var gh *gameHandlers.GameHandlers = gameHandlers.NewGameHandlers()
	http.HandleFunc("/games", gh.DistributeMethods)
	http.HandleFunc("/games/", gh.DistributeMethods)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
