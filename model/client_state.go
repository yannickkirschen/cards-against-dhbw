package model

//    JoinRequestAction (on error: InvalidCodeState)
//    LobbyState (Kind: WaitingLobbyState)
//      -> PlayerPlayingState && BossWaitingState   <-|
//      -> BossChoosingState                          |
//      -> BossHasChosenState                       --|
//      -> LobbyState (Kind: WinnerLobbyState)

<<<<<<< HEAD
type InvalidCodeState struct {
	Kind string `json:"kind"`
}

type JoinReqAction struct {
=======
type JoinRequestAction struct {
>>>>>>> e07cc93 (feat: add database and update states)
	Kind     string `json:"kind"`
	GameID   string `json:"gameID"`
	PlayerID string `json:"playerID"`
}

<<<<<<< HEAD
=======
type InvalidCodeState struct {
	Kind string `json:"kind"`
}

>>>>>>> e07cc93 (feat: add database and update states)
type LobbyState struct {
	Kind    string          `json:"kind"`
	Players []*PublicPlayer `json:"players"`
}

type StartGameAction struct {
	Kind string `json:"kind"`
}

type PlayerChoosingState struct {
	Kind       string          `json:"kind"`
	Players    []*PublicPlayer `json:"players"`
	BlackCard  *Card           `json:"blackCard"`
	WhiteCards []*Card         `json:"whiteCards"`
}

type CardChosenAction struct {
	Kind string `json:"kind"`
	Card string `json:"card"`
}

type BossWaitingState struct {
	Kind              string          `json:"kind"`
	Players           []*PublicPlayer `json:"players"`
	BlackCard         *Card           `json:"blackCard"`
	NumberPlayedCards int             `json:"numberPlayedCards"`
}

type BossHasChosenState struct {
	Kind        string                  `json:"kind"`
	Players     []*PublicPlayer         `json:"players"`
	BlackCard   *Card                   `json:"blackCard"`
	WinnerCard  *Card                   `json:"winnerCard"`
	PlayedCards map[*PublicPlayer]*Card `json:"playedCards"`
}

type LeaveGameAction struct {
	Kind string `json:"kind"`
}

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
