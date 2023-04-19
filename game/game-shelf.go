package game

import (
	"errors"
	"math/rand"
	"time"

	"github.com/yannickkirschen/cards-against-dhbw/data"
	"github.com/yannickkirschen/cards-against-dhbw/model"
	"github.com/yannickkirschen/cards-against-dhbw/utils"
)

var GlobalGameShelf *GameShelf

var ErrNoGame = errors.New("game does not exist")

type GameShelf struct {
	games map[string]*GamePlay
	r     *rand.Rand
}

func NewGameShelf() *GameShelf {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return &GameShelf{games: make(map[string]*GamePlay), r: r}
}

func (gs *GameShelf) CreateGame(name string) (string, string) {
	gameId := gs.newGameId()

	var blacks []*model.Card
	var whites []*model.Card
	data.ReadCards(blacks, whites)

	gamePlay := NewGamePlay(model.NewGame(gameId))
	player, _ := gamePlay.Game.GeneratePlayer(name)
	player.IsMod = true
	gamePlay.Game.Mod = player

	gs.games[gameId] = gamePlay
	return gameId, player.ID
}

func (gs *GameShelf) GamePlay(id string) (*GamePlay, error) {
	gp, exists := gs.games[id]
	if exists {
		return gp, nil
	}

	return nil, errors.New("game does not exist")
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

func (gs *GameShelf) JoinGame(gameId string, name string) (string, error) {
	game, ok := gs.games[gameId]
	if ok {
		player, err := game.Game.GeneratePlayer(name)
		if err != nil {
			return "", err
		}
		return player.ID, nil
	}
	return "", ErrNoGame
}
