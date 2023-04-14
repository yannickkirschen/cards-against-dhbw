package model

// The central game structure.
type Game struct {
	// The unique game code.
	Code string `json:"code"`

	// The player who is currently Master of Desaster.
	Mod *Player `json:"mod"`

	// The players in the game.
	Players []*Player `json:"players"`

	// The black cards in the game.
	// If a black card has been played, it is removed from the list.
	BlackCards []*Card `json:"blackCards"`

	// The white cards in the game.
	// If a white card has been played, it is removed from the list.
	WhiteCards []*Card `json:"whiteCards"`

	// The state of the game.
	State *State `json:"state"`
}
