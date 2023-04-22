package model

type JoinRequestAction struct {
	GameID   string `json:"gameID"`
	PlayerID string `json:"playerID"`
}

type InvalidState struct{}

type LobbyState struct {
	Players   []*PublicPlayer `json:"players"`
	GameReady bool            `json:"gameReady"`
}

type StartGameAction struct{}

type PlayerChoosingState struct {
	Players    []*PublicPlayer `json:"players"`
	BlackCard  *Card           `json:"blackCard"`
	WhiteCards []*Card         `json:"whiteCards"`
}

type CardChosenAction struct {
	Card string `json:"card"`
}

type BossHasChosenState struct {
	Players     []*PublicPlayer `json:"players"`
	BlackCard   *Card           `json:"blackCard"`
	Winner      string          `json:"winner"`
	WinnerCard  string          `json:"winnerCard"`
	PlayedCards []*Card         `json:"playedCards"`
}
