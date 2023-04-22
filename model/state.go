package model

// The state of the game.
type State struct {
	// The player who is currently choosing a white card.
	Boss *Player `json:"boss"`

	// The black card that is currently being played. If the black card is nil,
	// waiting to be started.
	BlackCard *Card `json:"blackCard"`

	// The white cards that have been played by the players. If the number of entries
	// is equal to the number of players, the round is over. Otherwise, players still
	// must choose a white card.
	PlayedCards map[*Player]*Card `json:"playedCards"`

	// The number of the current round. If the sum of all given points is less
	// than the round counter, a white card is being chosen by the Boss.
	Round int `json:"round"`
}

func NewState() *State {
	return &State{
		PlayedCards: make(map[*Player]*Card),
	}
}

type PublicPlayer struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IsMod  bool   `json:"isMod"`
	IsBoss bool   `json:"isBoss"`
	Points int    `json:"points"`
}

func (s *State) GetPlayedCards() []*Card {
	cards := make([]*Card, 0)
	for _, playedCard := range s.PlayedCards {
		cards = append(cards, playedCard)
	}
	return cards
}

func (s *State) WhoPlayed(cardId string) *Player {
	for player, card := range s.PlayedCards {
		if card.ID == cardId {
			return player
		}
	}

	return nil
}