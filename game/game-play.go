package game

import (
	"github.com/yannickkirschen/cards-against-dhbw/model"
)

type GamePlay struct {
	Game    *model.Game
	Senders map[string]func(title string, message ...any)
}

func NewGamePlay(game *model.Game) *GamePlay {
	return &GamePlay{
		Game:    game,
		Senders: make(map[string]func(title string, message ...any)),
	}
}

func (gp *GamePlay) UpdateState() {
	gp.Game.UpdatePublicPlayers()

	status := gp.Game.Status()
	if status == model.STATUS_BOSS_CHOOSING {
		gp.sendBossChoosingState()
	} else if status == model.STATUS_GAME_FINISHED {
		gp.sendLobbyState("game.finished")
	}
}

func (gp *GamePlay) AddSender(playerId string, sender func(string, ...any)) {
	gp.Senders[playerId] = sender
}

func (gp *GamePlay) ReceiveMessage(playerId string, title string, message any) {
	gp.Game.Mutex.Lock()
	defer gp.Game.Mutex.Unlock()

	switch title {
	case "game.join":
		if gp.Game.Status() == model.STATUS_GAME_LOBBY {
			gp.handleGameJoin()
			gp.UpdateState()
		}
	case "game.start":
		if gp.Game.Status() == model.STATUS_GAME_READY {
			gp.handleNewRound()
			gp.UpdateState()
		}
	case "player.card.chosen":
		if gp.Game.Status() == model.STATUS_PLAYER_CHOOSING {
			gp.handlePlayerCardChosenAction(playerId, message)
			gp.UpdateState()
		}
	case "boss.card.chosen":
		if gp.Game.Status() == model.STATUS_BOSS_CHOOSING {
			gp.handleBossCardChosenAction(playerId, message)
			gp.UpdateState()
		}
	case "boss.round.continue":
		if gp.Game.Status() == model.STATUS_ROUND_FINISHED {
			gp.handleNewRound()
			gp.Game.UpdatePublicPlayers()
		}
	default:
		gp.sendInvalidState(playerId)
	}
}

func (gp *GamePlay) handleGameJoin() {
	gp.sendLobbyState("game.joined")
}
func (gp *GamePlay) handleNewRound() {
	// Increment round
	gp.Game.State.Round++

	// New boss
	gp.Game.State.Boss = gp.Game.WhoIsNextBoss()

	// Clear played cards
	gp.Game.State.PlayedCards = make(map[*model.Player]*model.Card)

	// Choose black card and remove it from game
	if len(gp.Game.BlackCards) > 0 {
		gp.Game.State.BlackCard = gp.Game.ChooseBlackCard()
	}

	// Fill up white cards of all players and remove those from game
	for _, player := range gp.Game.Players {
		for len(player.Cards) < 10 {
			player.Cards = append(player.Cards, gp.Game.ChooseWhiteCard())
		}
	}

	gp.sendPlayerChoosingState()
}

func (gp *GamePlay) sendPlayerChoosingState() {
	for _, player := range gp.Game.Players {
		state := &model.PlayerChoosingState{
			Players:    gp.Game.PublicPlayers,
			BlackCard:  gp.Game.State.BlackCard,
			WhiteCards: player.Cards[:],
		}

		gp.Senders[player.ID]("player.choosing", state)
	}
}

func (gp *GamePlay) handlePlayerCardChosenAction(playerId string, message any) {
	action, ok := message.(model.CardChosenAction)
	player := gp.Game.FindPlayer(playerId)
	if ok && player == nil {
		gp.Game.State.PlayedCards[player] = gp.Game.FindCard(action.Card)
	}
}

func (gp *GamePlay) sendBossChoosingState() {
	state := &model.PlayerChoosingState{
		Players:    gp.Game.PublicPlayers,
		BlackCard:  gp.Game.State.BlackCard,
		WhiteCards: gp.Game.State.GetPlayedCards(),
	}

	for _, sender := range gp.Senders {
		sender("boss.choosing", state)
	}
}

func (gp *GamePlay) handleBossCardChosenAction(playerId string, message any) {
	action, ok := message.(model.CardChosenAction)
	if ok {
		winnerPlayer := gp.Game.State.WhoPlayed(action.Card)
		winnerCard := gp.Game.FindCard(action.Card)
		winnerPlayer.Points++
		gp.sendBossHasChosenState(winnerPlayer, winnerCard)
	}
}

func (gp *GamePlay) sendBossHasChosenState(winner *model.Player, winnerCard *model.Card) {
	state := &model.BossHasChosenState{
		Players:     gp.Game.PublicPlayers,
		BlackCard:   gp.Game.State.BlackCard,
		Winner:      winner.ID,
		WinnerCard:  winnerCard.ID,
		PlayedCards: gp.Game.State.GetPlayedCards(),
	}

	for _, sender := range gp.Senders {
		sender("boss.chosen", state)
	}
}

func (gp *GamePlay) sendLobbyState(title string) {
	state := &model.LobbyState{
		Players: gp.Game.PublicPlayers,
	}

	for _, sender := range gp.Senders {
		sender(title, state)
	}
}

func (gp *GamePlay) sendInvalidState(playerId string) {
	gp.Senders[playerId]("invalid", nil)
}
