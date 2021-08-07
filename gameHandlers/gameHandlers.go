package gameHandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	games "github.com/LYSingD/go-games-rest-api/games"
	uuid "github.com/google/uuid"
)

type Games []games.Game
type GameHandlers struct {
	Store map[string]games.Game // store = Game{}
	sync.Mutex
}

func NewGameHandlers() *GameHandlers {
	return &GameHandlers{
		Store: games.GameList,
	}
}

func (gh *GameHandlers) DistributeGamesMethods(w http.ResponseWriter, r *http.Request) {
	var methodsDistribution = map[string]interface{}{
		"GET":  gh.getGames,
		"POST": gh.postGame,
	}
	var requestMethod string = r.Method
	methodFunc, hasMethod := methodsDistribution[requestMethod]
	if !hasMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
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
		return
	}

	var newGame games.Game
	unmarshalErr := json.Unmarshal(bodyBytes, &newGame)
	if unmarshalErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(unmarshalErr.Error()))
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
	defer gh.Unlock()

	jsonBytes, err := json.Marshal(games) // write into json
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (gh *GameHandlers) DistributeGamesMethodsWithId(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimSpace(r.URL.Path)
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = strings.TrimSuffix(urlPath, "/")

	parameters := strings.Split(urlPath, "/")
	requiredNumOfParameters := 2
	hasNotEnoughParameters := requiredNumOfParameters > len(parameters)
	if hasNotEnoughParameters {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Insufficient Parameters\n"))
		return
	}

	var methodsDistributionWithId = map[string]interface{}{
		"GET": gh.getGameById,
		"PUT": gh.updateGameById,
	}
	var requestMethod string = r.Method
	methodFunc, hasMethod := methodsDistributionWithId[requestMethod]
	if !hasMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
		return
	}
	methodFunc.(func(http.ResponseWriter, *http.Request, string))(w, r, parameters[len(parameters)-1])

}

func (gh *GameHandlers) getGameById(w http.ResponseWriter, r *http.Request, id string) {
	gh.Lock()
	game, hasGame := gh.Store[id]
	defer gh.Unlock()

	if !hasGame {
		w.WriteHeader(http.StatusNotFound)
		responseMessage := fmt.Sprintf("Could not found game with id: %s\n", id)
		w.Write([]byte(responseMessage))
		return
	}

	jsonBytes, err := json.Marshal(game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// func (gh *GameHandlers) updateGameById(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Updating....")
// }
