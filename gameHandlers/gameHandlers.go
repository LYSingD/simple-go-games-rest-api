package gameHandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	games "github.com/LYSingD/go-games-rest-api/classes/games"
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
	var requestMethod string = r.Method
	switch requestMethod {
	case "GET":
		gh.getGames(w)
		return
	case "POST":
		gh.postGame(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
		return
	}
}

func (gh *GameHandlers) getGames(w http.ResponseWriter) {
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
	var requestMethod string = r.Method
	var urlGameId = parameters[len(parameters)-1]
	switch requestMethod {
	case "GET":
		gh.getGameById(w, r, urlGameId)
		return
	case "PUT":
		gh.updateGameById(w, r, urlGameId)
		return
	case "DELETE":
		gh.deleteGameById(w, r, urlGameId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
		return
	}
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

func (gh *GameHandlers) updateGameById(w http.ResponseWriter, r *http.Request, gameId string) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	gh.Lock()
	game, hasGame := gh.Store[gameId]
	gh.Unlock()
	if !hasGame {
		w.WriteHeader(http.StatusNotFound)
		responseMessage := fmt.Sprintf("Could not found game with id: %s\n", gameId)
		w.Write([]byte(responseMessage))
		return
	}

	unmarshalErr := json.Unmarshal(bodyBytes, &game)
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
	gh.Store[game.ID] = game
	delete(gh.Store, gameId)
	defer gh.Unlock()
}

func (gh *GameHandlers) deleteGameById(w http.ResponseWriter, r *http.Request, gameId string) {
	gh.Lock()
	_, isGameExisted := gh.Store[gameId]
	if !isGameExisted {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := fmt.Sprintf("Could not find the game with Id: %s\n", gameId)
		w.Write([]byte(errorMessage))
		return
	}
	delete(gh.Store, gameId)
	defer gh.Unlock()
}
