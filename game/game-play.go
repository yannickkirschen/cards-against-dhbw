package game

import (
	"log"

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
	log.Printf("Game %s (GamePlay): status is %d", gp.Game.Code, status)

	if status == model.STATUS_BOSS_CHOOSING {
		gp.sendBossChoosingState()
	} else if status == model.STATUS_GAME_FINISHED {
		gp.sendLobbyState("game.finished", true)
	}

	log.Printf("Game %s (GamePlay): updated state!", gp.Game.Code)
}

func (gp *GamePlay) AddSender(playerId string, sender func(string, ...any)) {
	gp.Senders[playerId] = sender
	log.Printf("Game %s (GamePlay): added sender for player %s!", gp.Game.Code, playerId)
}

func (gp *GamePlay) ReceiveMessage(playerId string, title string, message any) {
	gp.Game.Mutex.Lock()

	defer gp.UpdateState() // In order to react to changes in the state that may occur after an action has been processed
	defer gp.Game.Mutex.Unlock()

	log.Printf("Game %s (GamePlay): received message from player %s of type %s!", gp.Game.Code, playerId, title)
	gp.UpdateState()

	switch title {
	case "game.join":
		if gp.Game.Status() == model.STATUS_GAME_LOBBY || gp.Game.Status() == model.STATUS_GAME_READY {
			gp.handleGameJoin()
		}
	case "game.start":
		if gp.Game.Status() == model.STATUS_GAME_READY {
			gp.handleNewRound()
		}
	case "entity.card.chosen":
		if gp.Game.Status() == model.STATUS_PLAYER_CHOOSING {
			gp.handlePlayerCardChosenAction(playerId, message)
		} else if gp.Game.Status() == model.STATUS_BOSS_CHOOSING {
			gp.handleBossCardChosenAction(playerId, message)
		}
	case "mod.round.continue":
		if gp.Game.Status() == model.STATUS_ROUND_FINISHED {
			gp.handleNewRound()
		}
	default:
		log.Printf("Game %s (GamePlay): unrecognized message type %s from player %s!", gp.Game.Code, title, playerId)
		gp.sendInvalidState(playerId)
	}
}

func (gp *GamePlay) handleGameJoin() {
	gp.sendLobbyState("game.joined", gp.Game.Status() == model.STATUS_GAME_READY)
}

func (gp *GamePlay) handleNewRound() {
	if gp.Game.State.Boss == nil {
		log.Printf("Game %s (GamePlay): setting up new round. Previous round was %d (no boss yet)!", gp.Game.Code, gp.Game.State.Round)
	} else {
		log.Printf("Game %s (GamePlay): setting up new round. Previous round was %d (boss is %s)!", gp.Game.Code, gp.Game.State.Round, gp.Game.State.Boss.Name)
	}

	// Increment round
	gp.Game.State.Round++
	log.Printf("Game %s (GamePlay): new round is %d!", gp.Game.Code, gp.Game.State.Round)

	// New boss
	gp.Game.State.Boss = gp.Game.WhoIsNextBoss()
	log.Printf("Game %s (GamePlay): new boss is %s!", gp.Game.Code, gp.Game.State.Boss.Name)

	// Clear played cards
	gp.Game.State.PlayedCards = make(map[*model.Player]*model.Card)

	// Choose black card and remove it from game
	// TODO: handle error when no black cards are left
	if len(gp.Game.BlackCards) > 0 {
		gp.Game.State.BlackCard = gp.Game.ChooseBlackCard()
		log.Printf("Game %s (GamePlay): new black card is %s!", gp.Game.Code, gp.Game.State.BlackCard.Text)
	}

	// Fill up white cards of all players and remove those from game
	for _, player := range gp.Game.Players {
		for len(player.Cards) < 10 {
			c := gp.Game.ChooseWhiteCard()
			player.Cards = append(player.Cards, c)
		}
	}

	log.Printf("Game %s (GamePlay): all players have ten white cards again!", gp.Game.Code)
	log.Printf("Game %s (GamePlay): created new round!", gp.Game.Code)

	gp.UpdateState()
	gp.sendPlayerChoosingState()
}

func (gp *GamePlay) sendPlayerChoosingState() {
	for _, player := range gp.Game.Players {
		state := &model.PlayerChoosingState{
			Players:    gp.Game.PublicPlayers,
			BlackCard:  gp.Game.State.BlackCard,
			WhiteCards: player.Cards,
		}

		gp.Senders[player.ID]("player.choosing", state)
	}

	log.Printf("Game %s (GamePlay): sent message of type 'player.choosing' to all players!", gp.Game.Code)
}

func (gp *GamePlay) handlePlayerCardChosenAction(playerId string, message any) {
	action, ok := message.(model.CardChosenAction)
	player := gp.Game.FindPlayer(playerId)

	if ok && player != nil {
		card := gp.Game.FindCardEverywhere(action.Card, player)
		player.RemoveCard(card.ID)
		gp.Game.State.PlayedCards[player] = card
		gp.sendPlayerChoosingState()
		log.Printf("Game %s (GamePlay): player %s played card %s!", gp.Game.Code, player.Name, card.Text)
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
	log.Printf("Game %s (GamePlay): sent message of type 'boss.choosing' to all players!", gp.Game.Code)
}

func (gp *GamePlay) handleBossCardChosenAction(playerId string, message any) {
	action, ok := message.(model.CardChosenAction) // TODO: check if player is boss (or is this done in the socket code?)
	if ok {
		winnerPlayer := gp.Game.State.WhoPlayed(action.Card)
		winnerCard := gp.Game.FindCardEverywhere(action.Card, winnerPlayer)
		winnerPlayer.Points++

		gp.UpdateState()
		gp.sendBossHasChosenState(winnerPlayer, winnerCard)
		log.Printf("Game %s (GamePlay): boss %s chose card %s! Winner is %s!", gp.Game.Code, playerId, action.Card, winnerPlayer.Name)
	} else {
		log.Printf("Game %s (GamePlay): invalid message of type 'entity.card.chosen' from player %s!", gp.Game.Code, playerId)
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
	log.Printf("Game %s (GamePlay): sent message of type 'boss.chosen' with new game state to all players!", gp.Game.Code)
}

func (gp *GamePlay) sendLobbyState(title string, gameReady bool) {
	state := &model.LobbyState{
		Players:   gp.Game.PublicPlayers,
		GameReady: gameReady,
	}

	for _, sender := range gp.Senders {
		sender(title, state)
	}
	log.Printf("Game %s (GamePlay): sent message of type '%s' to all players!", gp.Game.Code, title)
}

func (gp *GamePlay) sendInvalidState(playerId string) {
	gp.Senders[playerId]("invalid", "")
	log.Printf("Game %s (GamePlay): sent message of type 'invalid' to player %s!", gp.Game.Code, playerId)
}
