package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yannickkirschen/cards-against-dhbw/utils"
)

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to create new game (endpoint: /v1/new, controller: NewGameHandler, method: %s))", r.Method)

	switch r.Method {
	case http.MethodGet:
		handleNewGameGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to join a game (endpoint: /v1/join, controller: JoinGameHandler, method: %s))", r.Method)

	switch r.Method {
	case http.MethodGet:
		handleJoinGameGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleNewGameGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	name := r.URL.Query().Get("name")
	if name == "" {
		log.Print("No game ID given in query parameters (name is empty). Returning 400 Bad Request.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no name given"))
		return
	}

	gameId, playerId, err := GlobalGameShelf.CreateGame(name)
	if err != nil {
		log.Printf("Error creating game: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error creating game: %s", err.Error())))
		return
	}

	response := &response{
		GameId:   gameId,
		PlayerId: playerId,
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error converting game creation response to JSON: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error converting response to JSON: %s", err.Error())))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func handleJoinGameGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	gameId, err := utils.PathParameterFilter(r.URL.Path, "/v1/join/")
	if err != nil {
		log.Printf("Error parsing game ID from URL: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		log.Print("No game ID given in query parameters (name is empty). Returning 400 Bad Request.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no name given"))
		return
	}

	playerId, err := GlobalGameShelf.JoinGame(gameId, name)
	switch err {
	case ErrNotFound:
		log.Printf("Player %s cannot join game %s as it does not exist. Error was: %s", name, gameId, err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Game not found"))
	case nil:
		response := &response{
			GameId:   gameId,
			PlayerId: playerId,
		}

		b, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error converting game joined response to JSON: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error converting response to JSON: %s", err.Error())))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

type response struct {
	GameId   string `json:"gameId"`
	PlayerId string `json:"playerId"`
}
