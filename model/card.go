package model

const (
	// Represents a black card.
	BLACK = 0

	// Represents a white card.
	WHITE = 1
)

// A card in the game.
// If it is a black card, it contains one blank represented by an underscore.
// If the card is a black card and does not contain an underscore, it is
// assumed to be a question card that also accepts one white card per player.
type Card struct {
	// Hash of the text.
	ID string `json:"id"`

	// The text of the card.
	Text string `json:"text"`

	// The type of the card.
	Type int `json:"type"`
}
