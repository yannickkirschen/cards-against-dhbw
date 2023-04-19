package game

import (
	"errors"
	"math/rand"
	"time"

	"github.com/yannickkirschen/cards-against-dhbw/model"
	"github.com/yannickkirschen/cards-against-dhbw/utils"
)

var GlobalGameShelf *GameShelf

var ErrNoGame = errors.New("game does not exist")

type GameShelf struct {
	Games map[string]*GamePlay
	r     *rand.Rand
}

func NewGameShelf() *GameShelf {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return &GameShelf{Games: make(map[string]*GamePlay), r: r}
}

func (gs *GameShelf) CreateGame(name string) (string, string) {
	gameId := gs.newGameId()
	gamePlay := NewGamePlay(model.NewGame(gameId))
	player, _ := gamePlay.Game.GeneratePlayer(name)
	player.IsMod = true
	gamePlay.Game.Mod = player

	//TODO: add black and white cards to gameplay.Game
	gs.Games[gameId] = gamePlay
	return gameId, player.ID
}

func (gs *GameShelf) newGameId() string {
	gameId := utils.RandString(gs.r, 4)
	for {
		_, exists := gs.Games[gameId]

		if exists {
			gameId = utils.RandString(gs.r, 4)
		} else {
			return gameId
		}
	}
}

func (gs *GameShelf) JoinGame(gameId string, name string) (string, error) {
	game, ok := gs.Games[gameId]
	if ok {
		player, err := game.Game.GeneratePlayer(name)
		if err != nil {
			return "", err
		}
		return player.ID, nil
	}
	return "", ErrNoGame
}
