package model

import (
	"encoding/base64"
	"errors"
	"sync"
)

const (
	STATUS_GAME_LOBBY = iota
	STATUS_GAME_READY
	STATUS_PLAYER_CHOOSING
	STATUS_BOSS_CHOOSING
	STATUS_ROUND_FINISHED
	STATUS_GAME_FINISHED
)

// The central game structure.
type Game struct {
	// The mutex used to synchronize access to the game.
	Mutex sync.Mutex

	// The unique game code.
	Code string `json:"code"`

	// The player who is currently Master of Disaster.
	Mod *Player `json:"mod"`

	// The players in the game.
	Players []*Player `json:"players"`

	// The players' public representation.
	PublicPlayers []*PublicPlayer `json:"publicPlayers"`

	// The black cards in the game.
	// If a black card has been played, it is removed from the list.
	BlackCards []*Card `json:"blackCards"`

	// The white cards in the game.
	// If a white card has been played, it is removed from the list.
	WhiteCards []*Card `json:"whiteCards"`

	// The state of the game.
	State *State `json:"state"`
}

func (g *Game) Status() int {
	// Status: Game is waiting for players to join.
	//  - There are less than two players in the game.
	if g.State.Round == 0 && len(g.Players) < 2 {
		return STATUS_GAME_LOBBY
	}

	// Status: Game is ready to be started.
	//  - Round counter is at zero.
	//  - There are at least two players in the game.
	if g.State.Round == 0 && len(g.Players) > 1 {
		return STATUS_GAME_READY
	}

	// Status: Players are choosing their white cards; Boss is waiting.
	//  - The sum of all points is n, while the current round is n+1.
	//  - Not all players have played a card.
	if g.State.Round-1 == g.SumOfPoints() && !g.State.AllCardsPlayed() {
		return STATUS_PLAYER_CHOOSING
	}

	// Status: Players are waiting for boss to choose the best white card.
	//  - The sum of all points is n, while the current round is n+1.
	//  - All players have played a card.
	if g.State.Round-1 == g.SumOfPoints() && g.State.AllCardsPlayed() {
		return STATUS_BOSS_CHOOSING
	}

	// Status: Boss has chosen a card, point was given, round is over.
	//  - The sum of all points is equal to the current round.
	//  - No player has 10 points.
	if g.State.Round == g.SumOfPoints() && !g.WeHaveAWinner() {
		return STATUS_ROUND_FINISHED
	}

	// Status: Boss has chosen a card, point was given, game is over.
	//  - The sum of all points is equal to the current round.
	//  - One player has 10 points.
	if g.State.Round == g.SumOfPoints() && g.WeHaveAWinner() {
		return STATUS_GAME_FINISHED
	}

	// Status: invalid
	return -1
}

func (g *Game) SumOfPoints() int {
	var sum = 0
	for _, player := range g.Players {
		sum += player.Points
	}
	return sum
}

func (g *Game) FindCard(id string) *Card {
	for _, card := range g.BlackCards {
		if card.ID == id {
			return card
		}
	}

	for _, card := range g.WhiteCards {
		if card.ID == id {
			return card
		}
	}

	return nil
}

func (g *Game) GeneratePlayer(name string) (*Player, error) {
	if g.PlayerNameExists(name) {
		return nil, errors.New("player exists")
	}
	player := &Player{
		ID:   base64.StdEncoding.EncodeToString([]byte(name)),
		Name: name,
	}
	return player, nil
}

func (g *Game) PlayerNameExists(name string) bool {
	for _, player := range g.Players {
		if player.Name == name {
			return true
		}
	}
	return false
}

func (g *Game) UpdatePublicPlayers() {
	g.PublicPlayers = nil
	for _, player := range g.Players {
		g.PublicPlayers = append(g.PublicPlayers, &PublicPlayer{
			Name:   player.Name,
			IsMod:  player.IsMod,
			IsBoss: (g.State.Boss == player),
			Points: player.Points,
		})
	}
}

func (g *Game) FindPlayer(id string) *Player {
	for _, player := range g.Players {
		if player.ID == id {
			return player
		}
	}

	return nil
}

func (g *Game) WeHaveAWinner() bool {
	for _, player := range g.Players {
		if player.Points == 10 {
			return true
		}
	}

	return false
}

func (g *Game) ChooseBlackCard() *Card {
	card := g.BlackCards[0]
	g.BlackCards = g.BlackCards[1:]
	return card
}

func (g *Game) ChooseWhiteCard() *Card {
	card := g.WhiteCards[0]
	g.WhiteCards = g.WhiteCards[1:]
	return card
}

func (g *Game) WhoIsNextBoss() *Player {
	highestIndex := len(g.Players) - 1
	for index, player := range g.Players {
		if g.State.Boss == player {
			return g.Players[(index+1)%highestIndex]
		}
	}

	return g.Players[0]
}
