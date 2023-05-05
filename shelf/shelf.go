package shelf

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/yannickkirschen/cards-against-dhbw/data"
	"github.com/yannickkirschen/cards-against-dhbw/game"
	"github.com/yannickkirschen/cards-against-dhbw/play"
	"github.com/yannickkirschen/cards-against-dhbw/utils"
)

var GlobalShelf *GameShelf
var ErrNotFound = errors.New("game does not exist")

type GameShelf struct {
	games map[string]*play.Play
	r     *rand.Rand
}

func New() *GameShelf {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return &GameShelf{
		games: make(map[string]*play.Play),
		r:     r,
	}
}

func (gs *GameShelf) CreateGame(name string) (string, string, error) {
	log.Printf("Player %s wants to create a new game!", name)

	gameId := gs.newGameId()
	log.Printf("Created a new game ID for the game player %s wants to create. ID: %s", name, gameId)

	blacks, whites, err := data.ReadCards()
	if err != nil {
		return "", "", err
	}

	p := play.New(game.New(gameId, blacks, whites))
	p.DeleteCallback = gs.DeleteCallback
	player, _ := p.Game.CreatePlayer(name)
	p.Game.Mod = player.Name

	gs.games[gameId] = p

	return gameId, player.Name, data.Update(p.Game)
}

func (gs *GameShelf) Play(id string) (*play.Play, error) {
	gp, exists := gs.games[id]
	if exists {
		return gp, nil
	} else {
		game, err := data.FindGame(id)
		if err != nil {
			return nil, ErrNotFound
		}

		gs.games[id] = play.New(game)
		gs.games[id].DeleteCallback = gs.DeleteCallback
		return gs.games[id], nil
	}
}

func (gs *GameShelf) JoinGame(gameId string, name string) (string, error) {
	log.Printf("Player %s wants to join game %s!", name, gameId)

	gamePlay, ok := gs.games[gameId]
	if ok {
		player, err := gamePlay.Game.CreatePlayer(name)
		if err != nil {
			return "", err
		}

		return player.Name, nil
	} else {
		game, err := data.FindGame(gameId)
		if err != nil {
			return "", ErrNotFound
		}

		gs.games[gameId] = play.New(game)
		gs.games[gameId].DeleteCallback = gs.DeleteCallback

		if game.FindPlayer(name) == nil {
			return gs.JoinGame(gameId, name)
		} else {
			return name, nil
		}
	}
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

func (gs *GameShelf) DeleteCallback(gameCode string) {
	log.Printf("Deleting game %s", gameCode)

	delete(gs.games, gameCode)

	err := data.Delete(gameCode)
	if err != nil {
		log.Printf("Error deleting game %s: %s", gameCode, err.Error())
	}
}
