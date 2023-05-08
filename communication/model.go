package communication

import (
	"github.com/yannickkirschen/cards-against-dhbw/card"
	"github.com/yannickkirschen/cards-against-dhbw/player"
)

type JoinRequestAction struct {
	GameCode   string `json:"gameCode"`
	PlayerName string `json:"playerName"`
}

type LobbyState struct {
	Players   []*player.PublicPlayer `json:"players"`
	GameReady bool                   `json:"gameReady"`
}

type StartGameAction struct{}

type PlayerChoosingState struct {
	Players    []*player.PublicPlayer `json:"players"`
	BlackCard  *card.Card             `json:"blackCard"`
	WhiteCards []*card.Card           `json:"whiteCards"`
}

type CardChosenAction struct {
	CardId string `json:"cardId"`
}

type BossHasChosenState struct {
	Players     []*player.PublicPlayer `json:"players"`
	BlackCard   *card.Card             `json:"blackCard"`
	Winner      string                 `json:"winner"`
	WinnerCard  string                 `json:"winnerCard"`
	PlayedCards []*card.Card           `json:"playedCards"`
}

type PlayerKickAction struct {
	PlayerName string `json:"playerName"`
}

type ApplicationError struct {
	Label   string `json:"label"`
	Payload any    `json:"payload"`
}
