package round

import (
	"github.com/yannickkirschen/cards-against-dhbw/card"
	"github.com/yannickkirschen/cards-against-dhbw/player"
)

type Round struct {
	// The player who is currently choosing a white card.
	Boss string `json:"boss"`

	// The black card that is currently being played. If the black card is nil,
	// waiting to be started.
	BlackCard *card.Card `json:"blackCard"`

	// The white cards that have been given to the players.
	WhiteCards map[string][]*card.Card `json:"whiteCards"`

	// The white cards that have been played by the players. If the number of entries
	// is equal to the number of players, the round is over. Otherwise, players still
	// must choose a white card.
	PlayedCards map[string]*card.Card `json:"playedCards"`

	// If the round is over, this field contains the winner of the round.
	Winner string `json:"winner"`

	// The number of the current round. If the sum of all given points is less
	// than the round counter, a white card is being chosen by the Boss.
	Counter int `json:"counter"`
}

func New() *Round {
	return &Round{
		WhiteCards:  make(map[string][]*card.Card),
		PlayedCards: make(map[string]*card.Card),
	}
}

func (r *Round) RemoveAllForPlayer(player *player.Player) {
	delete(r.WhiteCards, player.Name)
	delete(r.PlayedCards, player.Name)
}

// Find the card with the given ID in the white cards of the given player.
func (r *Round) FindCardFor(player *player.Player, id string) *card.Card {
	for _, card := range r.WhiteCards[player.Name] {
		if card.Id == id {
			return card
		}
	}

	return nil
}

// Remove the given card from the white cards of the given player.
func (r *Round) RemoveCardFor(player *player.Player, card *card.Card) {
	for i, c := range r.WhiteCards[player.Name] {
		if c == card {
			r.WhiteCards[player.Name] = append(r.WhiteCards[player.Name][:i], r.WhiteCards[player.Name][i+1:]...)
			return
		}
	}
}

// Returns the player who played the card with the given ID.
func (r *Round) WhoPlayed(id string) string {
	for player, card := range r.PlayedCards {
		if card.Id == id {
			return player
		}
	}

	return ""
}

// Returns a list of all played cards.
func (r *Round) FlatPlayedCards() []*card.Card {
	cards := make([]*card.Card, 0)
	for _, playedCard := range r.PlayedCards {
		cards = append(cards, playedCard)
	}
	return cards
}
