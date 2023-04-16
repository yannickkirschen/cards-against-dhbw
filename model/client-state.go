package model

//    JoinRequestAction (on error: InvalidCodeState)
//    LobbyState (Kind: WaitingLobbyState)
//      -> PlayerPlayingState && BossWaitingState   <-|
//      -> BossChoosingState                          |
//      -> BossHasChosenState                       --|
//      -> LobbyState (Kind: WinnerLobbyState)

type JoinRequestAction struct {
	GameID   string `json:"gameID"`
	PlayerID string `json:"playerID"`
}

type InvalidState struct{}

type LobbyState struct {
	Players []*PublicPlayer `json:"players"`
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

type LeaveGameAction struct{}

// MOD
// Visit example.com
//     Click create game
//     Enter name
// Call example.com/new
//     Game is created
//     Response is game code and player ID
// Redirect example.com/<game.id> with player ID in react state

// Other players
// Call example.com/join/<game.id>
//     Enter name
//     Click join
// Redirect example.com/<game.id> with player ID in react state

// All players
// example.com/<game.id> with player ID in react state
//     Internal call: socket example.com:1234
//         Send: game ID and player ID
//         Status: while listening loop
