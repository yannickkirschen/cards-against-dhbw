package play

import (
	"log"

	"github.com/yannickkirschen/cards-against-dhbw/card"
	"github.com/yannickkirschen/cards-against-dhbw/communication"
	"github.com/yannickkirschen/cards-against-dhbw/game"
	"github.com/yannickkirschen/cards-against-dhbw/player"
)

type Play struct {
	Game    *game.Game
	senders map[*player.Player]communication.Sender
}

func New(game *game.Game) *Play {
	return &Play{
		Game:    game,
		senders: make(map[*player.Player]communication.Sender),
	}
}

func (p *Play) AddSender(playerName string, sender communication.Sender) {
	p.senders[p.Game.FindPlayer(playerName)] = sender
}

func (p *Play) Receive(player string, action string, message any) {
	p.Game.Mutex.Lock()
	p.Game.UpdateState()

	defer p.Game.UpdateState()
	defer p.Game.Mutex.Unlock()

	log.Printf("Received message from player %s of type %s for game %s!", player, action, p.Game.Code)

	if !p.Game.StateAllows(action) {
		// TODO: send ACTION_FORBIDDEN
	}

	switch action {
	case game.ACTION_GAME_JOIN:
		p.sendLobbyState(game.STATE_GAME_LOBBY, p.Game.State == game.STATE_GAME_READY)
	case game.ACTION_GAME_START:
		fallthrough
	case game.ACTION_ROUND_CONTINUE:
		p.Game.StartNewRound()
		p.sendPlayersChoosingState()
	case game.ACTION_GAME_LEAVE:
		p.handlePlayerLeave(player)
	case game.ACTION_CARD_CHOSEN:
		p.handlePlayerChosenAction(player, message)
	default:
		p.sendInvalidState(player)
	}

	p.Game.UpdateState()
	if p.Game.State == game.STATE_BOSS_CHOOSING {
		p.sendBossChoosingState()
	} else if p.Game.State == game.STATE_GAME_FINISHED {
		p.sendLobbyState(game.STATE_GAME_FINISHED, true)
	}
}

func (p *Play) sendPlayersChoosingState() {
	for _, player := range p.Game.Players {
		state := &communication.PlayerChoosingState{
			Players:    p.Game.GeneratePublicPlayers(),
			BlackCard:  p.Game.CurrentRound.BlackCard,
			WhiteCards: p.Game.CurrentRound.WhiteCards[player],
		}

		p.senders[player].Send(game.STATE_PLAYERS_CHOOSING, state)
	}

	log.Printf("Sent state '%s' to all players for game %s!", game.STATE_PLAYERS_CHOOSING, p.Game.Code)
}

func (p *Play) sendBossChoosingState() {
	state := &communication.PlayerChoosingState{
		Players:    p.Game.GeneratePublicPlayers(),
		BlackCard:  p.Game.CurrentRound.BlackCard,
		WhiteCards: p.Game.CurrentRound.FlatPlayedCards(),
	}

	for _, sender := range p.senders {
		sender.Send("boss.choosing", state)
	}
}

func (p *Play) handlePlayerChosenAction(playerName string, message any) {
	action, ok := message.(communication.CardChosenAction)
	player := p.Game.FindPlayer(playerName)

	if !ok || player == nil {
		return
	}

	if p.Game.State == game.STATE_PLAYERS_CHOOSING {
		card := p.Game.CurrentRound.FindCardFor(player, action.Id)
		p.Game.CurrentRound.PlayedCards[player] = card
		p.Game.CurrentRound.RemoveCardFor(player, card)
		p.sendPlayersChoosingState()
	} else if p.Game.State == game.STATE_BOSS_CHOOSING {
		winnerPlayer := p.Game.CurrentRound.WhoPlayed(action.Id)
		winnerPlayer.Points++
		winnerCard := p.Game.CurrentRound.PlayedCards[winnerPlayer]
		p.Game.CurrentRound.Winner = p.Game.CurrentRound.WhoPlayed(action.Id)

		log.Printf("Player %s won round %d with card %s!", winnerPlayer.Name, p.Game.CurrentRound.Counter, winnerCard.Text)

		p.Game.UpdateState()
		p.sendBossHasChosenState(winnerPlayer, winnerCard)
	}
}

func (p *Play) sendBossHasChosenState(winner *player.Player, winnerCard *card.Card) {
	state := &communication.BossHasChosenState{
		Players:     p.Game.GeneratePublicPlayers(),
		BlackCard:   p.Game.CurrentRound.BlackCard,
		Winner:      winner.Name,
		WinnerCard:  winnerCard.Text,
		PlayedCards: p.Game.CurrentRound.FlatPlayedCards(),
	}

	for _, sender := range p.senders {
		sender.Send(game.STATE_ROUND_FINISHED, state)
	}

	log.Printf("Sent state '%s' to all players for game %s!", game.STATE_ROUND_FINISHED, p.Game.Code)
}

func (p *Play) sendLobbyState(stateName string, gameReady bool) {
	state := &communication.LobbyState{
		Players:   p.Game.GeneratePublicPlayers(),
		GameReady: gameReady,
	}

	for _, sender := range p.senders {
		sender.Send(stateName, state)
	}

	log.Printf("Sent lobby state to all players for game %s (state was '%s')!", p.Game.Code, stateName)
}

func (p *Play) handlePlayerLeave(playerName string) {
	player := p.Game.FindPlayer(playerName)
	if player == nil {
		return
	}

	p.Game.RemovePlayer(player)

	if len(p.Game.Players) < 2 {
		p.Game.State = game.STATE_GAME_LOBBY
		p.Game.CurrentRound = nil
		p.sendLobbyState(game.STATE_GAME_LOBBY, false)
		return
	}

	if p.Game.Mod == player {
		p.Game.Mod = p.Game.Players[0]
	}

	if p.Game.CurrentRound.Boss == player {
		p.Game.StartNewRound()
		p.sendPlayersChoosingState()
	} else {
		p.Game.CurrentRound.RemoveAllForPlayer(player)

		switch p.Game.State {
		case game.STATE_GAME_LOBBY:
			p.sendLobbyState(game.STATE_GAME_LOBBY, false)
		case game.STATE_PLAYERS_CHOOSING:
			p.sendPlayersChoosingState()
		case game.STATE_BOSS_CHOOSING:
			p.sendBossChoosingState()
		case game.STATE_ROUND_FINISHED:
			p.sendBossHasChosenState(p.Game.CurrentRound.Winner, p.Game.CurrentRound.PlayedCards[p.Game.CurrentRound.Winner])
		case game.STATE_GAME_FINISHED:
			p.sendLobbyState(game.STATE_GAME_FINISHED, true)
		default:
			p.sendInvalidState(playerName)
		}
	}
}

func (p *Play) sendInvalidState(player string) {
	p.senders[p.Game.FindPlayer(player)].Send(game.STATE_INVALID, "")
	log.Printf("Sent invalid state to player %s for game %s!", player, p.Game.Code)
}
