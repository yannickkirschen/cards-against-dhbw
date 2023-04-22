package card

const (
	BLACK = 0 // Represents a black card.
	WHITE = 1 // Represents a white card.
)

// A card in the game.
// If it is a black card, it contains one blank represented by an underscore.
// If the card is a black card and does not contain an underscore, it is
// assumed to be a question card that also accepts one white card per player.
type Card struct {
	Id   string `json:"id"`
	Text string `json:"text"`
	Type int    `json:"type"`
}
