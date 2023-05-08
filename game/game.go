package game

import (
	"errors"
	"log"
	"sync"

	"github.com/yannickkirschen/cards-against-dhbw/card"
	"github.com/yannickkirschen/cards-against-dhbw/player"
	"github.com/yannickkirschen/cards-against-dhbw/round"
)

const (
	ACTION_GAME_JOIN       = "game.join"          // Player wants to join the game
	ACTION_GAME_START      = "game.start"         // MOD wants to start the game
	ACTION_GAME_LEAVE      = "game.leave"         // Player wants to leave the game
	ACTION_CARD_CHOSEN     = "entity.card.chosen" // Player has chosen a card (black or white)
	ACTION_ROUND_CONTINUE  = "mod.round.continue" // MOD wants to continue the round
	ACTION_PLAYER_INACTIVE = "player.inactive"    // Player is inactive
	ACTION_PLAYER_KICK     = "player.kick"        // MOD wants to kick a player
	ACTION_FORBIDDEN       = "forbidden"          // Action is forbidden in the current state
	ACTION_INVALID         = "invalid"            // Invalid action

	STATE_GAME_LOBBY       = "game.lobby"      // Players are waiting in the lobby
	STATE_GAME_READY       = "game.ready"      // Game is ready to be started
	STATE_PLAYERS_CHOOSING = "player.choosing" // Players are choosing their white cards
	STATE_BOSS_CHOOSING    = "boss.choosing"   // Boss is choosing the best white card
	STATE_ROUND_FINISHED   = "boss.chosen"     // Round is over (boss has chosen a card)
	STATE_GAME_FINISHED    = "game.finished"   // Game is over (one player has 10 points)
	STATE_INVALID          = "invalid"         // Invalid state

	MAX_POINTS = 10
)

type Game struct {
	// The mutex used to synchronize access to the game.
	Mutex sync.Mutex

	// The unique game code.
	Code string `json:"code"`

	// The state of the game.
	State string `json:"current"`

	// The player who is currently Master of Disaster.
	Mod string `json:"mod"`

	// The players in the game.
	Players []*player.Player `json:"players"`

	// The black cards in the game.
	// If a black card has been played, it is removed from the list.
	BlackCards []*card.Card `json:"blackCards"`

	// The white cards in the game.
	// If a white card has been played, it is removed from the list.
	WhiteCards []*card.Card `json:"whiteCards"`

	// The currently played round. If nil, the game is not running.
	CurrentRound *round.Round `json:"currentRound"`
}

// Creates a new game with the given code.
// Players, black cards and white cards are initialized as empty slices.
// The state is set to STATE_GAME_LOBBY.
func New(code string, blacks []*card.Card, whites []*card.Card) *Game {
	return &Game{
		Code:         code,
		State:        STATE_GAME_LOBBY,
		Players:      make([]*player.Player, 0),
		BlackCards:   blacks,
		WhiteCards:   whites,
		CurrentRound: round.New(),
	}
}

// Updates the state of this game.
// After calling this method, `Game.State` will be set to the new state.
func (g *Game) UpdateState() {
	log.Printf("Updating state of game %s (old state was '%s').", g.Code, g.State)

	if g.CurrentRound == nil || (g.CurrentRound.Counter == 0 && len(g.Players) < 2) {
		g.State = STATE_GAME_LOBBY
	} else if g.CurrentRound.Counter == 0 && len(g.Players) > 1 {
		g.State = STATE_GAME_READY
	} else if g.CurrentRound.Counter-1 == g.SumOfPoints() && !g.AllCardsPlayed() {
		g.State = STATE_PLAYERS_CHOOSING
	} else if g.CurrentRound.Counter-1 == g.SumOfPoints() && g.AllCardsPlayed() {
		g.State = STATE_BOSS_CHOOSING
	} else if g.AllCardsPlayed() && !g.WeHaveAWinner() {
		g.State = STATE_ROUND_FINISHED
	} else if g.AllCardsPlayed() && g.WeHaveAWinner() {
		g.State = STATE_GAME_FINISHED
	} else {
		g.State = STATE_INVALID
	}

	log.Printf("New state of game %s is '%s'.", g.Code, g.State)
}

func (g *Game) StateAllows(action string) bool {
	switch action {
	case ACTION_GAME_JOIN:
		return g.State == STATE_GAME_LOBBY || g.State == STATE_GAME_READY
	case ACTION_GAME_START:
		return g.State == STATE_GAME_READY
	case ACTION_GAME_LEAVE:
		return true
	case ACTION_CARD_CHOSEN:
		return g.State == STATE_PLAYERS_CHOOSING || g.State == STATE_BOSS_CHOOSING
	case ACTION_ROUND_CONTINUE:
		return g.State == STATE_ROUND_FINISHED
	case ACTION_PLAYER_KICK:
		return true // players can always be kicked
	default:
		return false
	}
}

// Returns the sum of all points of all players in this game.
// The sum is always less than or equal to the current round counter.
func (g *Game) SumOfPoints() int {
	var sum = 0
	for _, player := range g.Players {
		sum += player.Points
	}
	return sum
}

// This handy function checks, if all players (except the boss) have played a card
// in the current round. If so, it returns true. Otherwise, it returns false.
func (g *Game) AllCardsPlayed() bool {
	return len(g.CurrentRound.PlayedCards) == len(g.Players)-1
}

// Checks, if one of the players reached the maximum amount of points.
func (g *Game) WeHaveAWinner() bool {
	for _, player := range g.Players {
		if player.Points == MAX_POINTS {
			return true
		}
	}

	return false
}

func (g *Game) CreatePlayer(name string) (*player.Player, error) {
	if g.PlayerNameExists(name) {
		return nil, errors.New("player exists")
	}

	p := player.New(name)
	g.Players = append(g.Players, p)
	g.CurrentRound.WhiteCards[p.Name] = make([]*card.Card, 0)

	return p, nil
}

func (g *Game) RemovePlayer(player *player.Player) {
	var index = -1
	for i, p := range g.Players {
		if p == player {
			index = i
			break
		}
	}

	g.Players = append(g.Players[:index], g.Players[index+1:]...)
	if g.Mod == player.Name {
		g.Mod = ""
	}
}

func (g *Game) PlayerNameExists(name string) bool {
	for _, player := range g.Players {
		if player.Name == name {
			return true
		}
	}
	return false
}

// Generates a list of public players for this game.
func (g *Game) GeneratePublicPlayers() []*player.PublicPlayer {
	var boss string = ""
	if g.CurrentRound != nil {
		boss = g.CurrentRound.Boss
	}

	var publicPlayers = make([]*player.PublicPlayer, len(g.Players))
	for i, player := range g.Players {
		publicPlayers[i] = player.ToPublic(g.Mod, boss)
	}

	return publicPlayers
}

// Starts a new round.
func (g *Game) StartNewRound() {
	log.Printf("Starting new round in game %s.", g.Code)

	round := round.New()
	round.Counter = g.CurrentRound.Counter + 1
	round.PlayedCards = make(map[string]*card.Card)
	round.Boss = g.whoIsNextBoss()
	round.BlackCard = g.chooseNextBlackCard()

	round.WhiteCards = g.CurrentRound.WhiteCards
	g.fillUpWhiteCards()

	g.CurrentRound = round
	log.Printf("New round in game %s is %d.", g.Code, round.Counter)
}

func (g *Game) whoIsNextBoss() string {
	for index, player := range g.Players {
		if g.CurrentRound.Boss == player.Name {
			return g.Players[(index+1)%len(g.Players)].Name
		}
	}

	return g.Players[0].Name
}

func (g *Game) fillUpWhiteCards() {
	for player, cards := range g.CurrentRound.WhiteCards {
		new := cards[:]
		for len(new) < 10 {
			c := g.chooseNextWhiteCard()
			new = append(new, c)
		}
		g.CurrentRound.WhiteCards[player] = new
	}
}

func (g *Game) chooseNextBlackCard() *card.Card {
	if len(g.BlackCards) == 0 {
		panic("No black cards left in deck.")
	}

	card := g.BlackCards[0]
	if len(g.BlackCards) > 0 {
		g.BlackCards = g.BlackCards[1:]
	}

	return card
}

func (g *Game) chooseNextWhiteCard() *card.Card {
	if len(g.WhiteCards) == 0 {
		panic("No white cards left in deck.")
	}

	card := g.WhiteCards[0]
	if len(g.WhiteCards) > 0 {
		g.WhiteCards = g.WhiteCards[1:]
	}

	return card
}

// Returns the player with the given name or nil, if no player with this name exists.
func (g *Game) FindPlayer(name string) *player.Player {
	for _, player := range g.Players {
		if player.Name == name {
			return player
		}
	}

	return nil
}
