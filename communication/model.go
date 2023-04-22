package communication

import (
	"github.com/yannickkirschen/cards-against-dhbw/card"
	"github.com/yannickkirschen/cards-against-dhbw/player"
)

type JoinRequestAction struct {
	GameID   string `json:"gameID"`
	PlayerID string `json:"playerID"`
}

type InvalidState struct{}

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
	Id string `json:"card"`
}

type BossHasChosenState struct {
	Players     []*player.PublicPlayer `json:"players"`
	BlackCard   *card.Card             `json:"blackCard"`
	Winner      string                 `json:"winner"`
	WinnerCard  string                 `json:"winnerCard"`
	PlayedCards []*card.Card           `json:"playedCards"`
}
