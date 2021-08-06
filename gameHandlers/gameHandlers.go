package gameHandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	games "github.com/LYSingD/go-games-rest-api/games"
	uuid "github.com/google/uuid"
)

type Games []games.Game

type GameHandlers struct {
	Store map[string]games.Game // store = Game{}
	sync.Mutex
}

func (gh *GameHandlers) DistributeMethods(w http.ResponseWriter, r *http.Request) {
	var methodsDistribution = map[string]interface{}{
		"GET":  gh.getGames,
		"POST": gh.postGame,
	}
	var requestMethod string = r.Method
	methodFunc, hasMethod := methodsDistribution[requestMethod]
	if !hasMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	methodFunc.(func(http.ResponseWriter, *http.Request))(w, r)

}

func (gh *GameHandlers) postGame(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Fatal(err)
		return
	}

	var newGame games.Game
	unmarshalErr := json.Unmarshal(bodyBytes, &newGame)
	if unmarshalErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		errorMessage := fmt.Sprintf("Need content-type 'application/json', but received '%s'", contentType)
		w.Write([]byte(errorMessage))
		return
	}

	gh.Lock()
	_, isGameExisted := gh.Store[newGame.ID]
	if !isGameExisted {
		newGame.ID = uuid.NewString()
		// Use current time as id instead of uuid
		// newGame.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	gh.Store[newGame.ID] = newGame
	defer gh.Unlock()

}

func (gh *GameHandlers) getGames(w http.ResponseWriter, r *http.Request) {
	var games Games

	gh.Lock()
	for _, game := range gh.Store {
		games = append(games, game)
	}
	gh.Unlock()

	jsonBytes, err := json.Marshal(games) // write into json
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Fatal(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func NewGameHandlers() *GameHandlers {
	return &GameHandlers{
		Store: map[string]games.Game{
			// "id1": {
			// 	ID:        "id1",
			// 	Name:      "Monster Hunter: World",
			// 	Developer: "Capcom",
			// 	Rating:    "T",
			// 	Genres:    []string{"Role-Playing", "Action RPG"},
			// },
		},
	}
}
