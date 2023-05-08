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

func (gs *GameShelf) CreateGame(playerName string) (string, string, error) {
	log.Printf("Player %s wants to create a new game!", playerName)

	gameCode := gs.newGameCode()
	log.Printf("Created a new game code for the game player %s wants to create. Code: %s", playerName, gameCode)

	blacks, whites, err := data.ReadCards()
	if err != nil {
		return "", "", err
	}

	p := play.New(game.New(gameCode, blacks, whites))
	p.DeleteCallback = gs.DeleteCallback
	player, _ := p.Game.CreatePlayer(playerName)
	p.Game.Mod = player.Name

	gs.games[gameCode] = p
	return gameCode, player.Name, data.Update(p.Game)
}

func (gs *GameShelf) Play(gameCode string) (*play.Play, error) {
	gp, exists := gs.games[gameCode]
	if exists {
		return gp, nil
	} else {
		game, err := data.FindGame(gameCode)
		if err != nil {
			return nil, ErrNotFound
		}

		gs.games[gameCode] = play.New(game)
		gs.games[gameCode].DeleteCallback = gs.DeleteCallback
		return gs.games[gameCode], nil
	}
}

func (gs *GameShelf) JoinGame(gameCode string, name string) (string, error) {
	log.Printf("Player %s wants to join game %s!", name, gameCode)

	gamePlay, ok := gs.games[gameCode]
	if ok {
		player, err := gamePlay.Game.CreatePlayer(name)
		if err != nil {
			return "", err
		}

		return player.Name, nil
	} else {
		game, err := data.FindGame(gameCode)
		if err != nil {
			return "", ErrNotFound
		}

		gs.games[gameCode] = play.New(game)
		gs.games[gameCode].DeleteCallback = gs.DeleteCallback

		if game.FindPlayer(name) == nil {
			return gs.JoinGame(gameCode, name)
		} else {
			return name, nil
		}
	}
}

func (gs *GameShelf) newGameCode() string {
	gameCode := utils.RandString(gs.r, 4)
	for {
		_, exists := gs.games[gameCode]

		if exists {
			gameCode = utils.RandString(gs.r, 4)
		} else {
			return gameCode
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
