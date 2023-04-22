package game

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/yannickkirschen/cards-against-dhbw/data"
	"github.com/yannickkirschen/cards-against-dhbw/model"
	"github.com/yannickkirschen/cards-against-dhbw/utils"
)

var GlobalGameShelf *GameShelf

var ErrNotFound = errors.New("game does not exist")

type GameShelf struct {
	games map[string]*GamePlay
	r     *rand.Rand
}

func NewGameShelf() *GameShelf {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return &GameShelf{games: make(map[string]*GamePlay), r: r}
}

func (gs *GameShelf) CreateGame(name string) (string, string, error) {
	log.Printf("Player %s wants to create a new game!", name)

	gameId := gs.newGameId()
	log.Printf("Created a new game ID for the game player %s wants to create. ID: %s", name, gameId)

	blacks, whites, err := data.ReadCards()
	if err != nil {
		return "", "", err
	}

	gamePlay := NewGamePlay(model.NewGame(gameId, blacks, whites))
	player, _ := gamePlay.Game.GeneratePlayer(name)
	player.IsMod = true

	gamePlay.Game.AddPlayer(player)
	gamePlay.Game.Mod = player

	gs.games[gameId] = gamePlay
	return gameId, player.ID, nil
}

func (gs *GameShelf) GamePlay(id string) (*GamePlay, error) {
	gp, exists := gs.games[id]
	if exists {
		return gp, nil
	}

	return nil, ErrNotFound
}

func (gs *GameShelf) JoinGame(gameId string, name string) (string, error) {
	log.Printf("Player %s wants to join game %s!", name, gameId)

	game, ok := gs.games[gameId]
	if ok {
		player, err := game.Game.GeneratePlayer(name)
		if err != nil {
			return "", err
		}

		game.Game.AddPlayer(player)
		game.Game.UpdatePublicPlayers()
		return player.ID, nil
	}

	return "", ErrNotFound
}

func (gs *GameShelf) newGameId() string {
	gameId := utils.RandString(gs.r, 4)
	for {
		_, exists := gs.games[gameId]

		if exists {
			gameId = utils.RandString(gs.r, 4)
		} else {
			return gameId
		}
	}
}
