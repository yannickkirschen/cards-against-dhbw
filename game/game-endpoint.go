package game

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yannickkirschen/cards-against-dhbw/utils"
)

func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleNewGameGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleJoinGameGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleNewGameGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no name given"))
		return
	}

	gameId, playerId := GlobalGameShelf.CreateGame(name)
	response := &response{
		GameId:   gameId,
		PlayerId: playerId,
	}

	b, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error converting response to JSON: %s", err.Error())))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func handleJoinGameGet(w http.ResponseWriter, r *http.Request) {
	gameId, err := utils.PathParameterFilter(r.URL.Path, "/v1/join/")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no name given"))
		return
	}

	playerId, err := GlobalGameShelf.JoinGame(gameId, name)
	switch err {
	case ErrNoGame:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Game not found"))
	case nil:
		response := &response{
			GameId:   gameId,
			PlayerId: playerId,
		}

		b, err := json.Marshal(response)
		if err != nil {
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
