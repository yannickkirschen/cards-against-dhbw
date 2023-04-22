package model

// A player in the game.
type Player struct {
	// Base64 of the name.
	ID string `json:"id"`

	// Unique name.
	Name string `json:"name"`

	// Whether the player is Master of Disaster or not.
	IsMod bool `json:"isMod"`

	// A player can hold ten white cards.
	Cards []*Card `json:"cards"`

	// The number of points the player has. For each won round there is one point.
	// If a player has 10 points. they win the game.
	Points int `json:"points"`
}

func NewPlayer(id string, name string, isMod bool) *Player {
	return &Player{
		ID:    id,
		Name:  name,
		IsMod: isMod,
		Cards: make([]*Card, 0),
	}
}
