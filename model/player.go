package model

import "net"

// A player in the game.
type Player struct {
	// Base64 of the name.
	ID string `json:"id"`

	// Unique name.
	Name string `json:"name"`

	// The socket connection of the player.
	Connection *net.Conn `json:"-"`

	// Whether the player is Master of Desaster or not.
	IsMod bool `json:"isMod"`

	// A player can hold ten white cards.
	Cards [10]*Card `json:"cards"`

	// The number of points the player has. For each won round there is one point.
	// If a player has 10 points. they win the game.
	Points int `json:"points"`
}

type PublicPlayer struct {
	Name   string `json:"name"`
	IsMod  bool   `json:"isMod"`
	IsBoss bool   `json:"isBoss"`
	Points int    `json:"points"`
}
