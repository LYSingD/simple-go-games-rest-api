package main

import (
	"log"
	"net/http"

	gameHandlers "github.com/LYSingD/go-games-rest-api/gameHandlers"
)

func main() {
	var gh *gameHandlers.GameHandlers = gameHandlers.NewGameHandlers()
	http.HandleFunc("/games", gh.DistributeGamesMethods)
	http.HandleFunc("/games/", gh.DistributeGamesMethodsWithId)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
