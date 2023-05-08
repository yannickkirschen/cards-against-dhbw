package play

import (
	"errors"
	"log"

	"github.com/yannickkirschen/cards-against-dhbw/card"
	"github.com/yannickkirschen/cards-against-dhbw/communication"
	"github.com/yannickkirschen/cards-against-dhbw/data"
	"github.com/yannickkirschen/cards-against-dhbw/err"
	"github.com/yannickkirschen/cards-against-dhbw/game"
)

type Play struct {
	Game           *game.Game
	DeleteCallback func(gameCode string)

	senders map[string]communication.Sender
}

func New(game *game.Game) *Play {
	return &Play{
		Game:    game,
		senders: make(map[string]communication.Sender),
	}
}

func (p *Play) AddSender(playerName string, sender communication.Sender) error {
	player := p.Game.FindPlayer(playerName)
	if player == nil {
		return errors.New("player not found")
	}

	p.senders[player.Name] = sender
	return nil
}

func (p *Play) Receive(playerName string, action string, message any) {
	p.Game.Mutex.Lock()
	p.Game.UpdateState()

	defer p.Game.Mutex.Unlock()

	log.Printf("Received message from player %s of type %s for game %s!", playerName, action, p.Game.Code)

	if !p.Game.StateAllows(action) {
		p.sendError(playerName, err.ACTION_FORBIDDEN, action)
		return
	}

	switch action {
	case game.ACTION_GAME_JOIN:
		p.Game.FindPlayer(playerName).Active = true
		log.Printf("Set player %s active!", playerName)
		p.sendCurrentState(playerName)
	case game.ACTION_GAME_START:
		fallthrough
	case game.ACTION_ROUND_CONTINUE:
		p.Game.StartNewRound()
		p.sendPlayersChoosingState()
	case game.ACTION_GAME_LEAVE:
		p.handlePlayerLeave(playerName)
		return
	case game.ACTION_CARD_CHOSEN:
		p.handlePlayerChosenAction(playerName, message)
	case game.ACTION_PLAYER_INACTIVE:
		p.Game.FindPlayer(playerName).Active = false
		log.Printf("Set player %s inactive!", playerName)
		p.sendCurrentState(playerName)
	case game.ACTION_PLAYER_KICK:
		p.handlePlayerKick(playerName, message)
		p.sendCurrentState(playerName)
	default:
		p.sendError(playerName, err.INVALID_STATE, p.Game.State)
	}

	p.Game.UpdateState()
	if p.Game.State == game.STATE_BOSS_CHOOSING {
		p.sendBossChoosingState()
	} else if p.Game.State == game.STATE_GAME_FINISHED {
		p.sendLobbyState(game.STATE_GAME_FINISHED, true)
	}

	p.Game.UpdateState()
	data.Update(p.Game)
}

func (p *Play) sendPlayersChoosingState() {
	for player, sender := range p.senders {
		state := &communication.PlayerChoosingState{
			Players:    p.Game.GeneratePublicPlayers(),
			BlackCard:  p.Game.CurrentRound.BlackCard,
			WhiteCards: p.Game.CurrentRound.WhiteCards[player],
		}

		sender.Send(game.STATE_PLAYERS_CHOOSING, state)
		log.Printf("Sent state '%s' to player %s for game %s!", game.STATE_PLAYERS_CHOOSING, player, p.Game.Code)
	}
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
		card := p.Game.CurrentRound.FindCardFor(player, action.CardId)
		if card == nil {
			p.sendError(playerName, err.CARD_NOT_FOUND, action.CardId)
			return
		}

		p.Game.CurrentRound.PlayedCards[player.Name] = card
		p.Game.CurrentRound.RemoveCardFor(player, card)
		p.sendPlayersChoosingState()
	} else if p.Game.State == game.STATE_BOSS_CHOOSING {
		playerName := p.Game.CurrentRound.WhoPlayed(action.CardId)
		if playerName == "" {
			p.sendError(playerName, err.CARD_NOT_PLAYED, action.CardId)
			return
		}

		winnerPlayer := p.Game.FindPlayer(playerName)
		if winnerPlayer == nil {
			p.sendError(p.Game.CurrentRound.Boss, err.PLAYER_NOT_FOUND, playerName)
			return
		}

		winnerPlayer.Points++
		winnerCard := p.Game.CurrentRound.PlayedCards[winnerPlayer.Name]
		p.Game.CurrentRound.Winner = winnerPlayer.Name

		log.Printf("Player %s won round %d with card %s!", winnerPlayer.Name, p.Game.CurrentRound.Counter, winnerCard.Text)

		p.Game.UpdateState()
		p.sendBossHasChosenState(winnerPlayer.Name, winnerCard)
	}
}

func (p *Play) sendBossHasChosenState(winner string, winnerCard *card.Card) {
	state := &communication.BossHasChosenState{
		Players:     p.Game.GeneratePublicPlayers(),
		BlackCard:   p.Game.CurrentRound.BlackCard,
		Winner:      winner,
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
	log.Printf("Player %s left game %s!", player.Name, p.Game.Code)

	if len(p.Game.Players) == 0 {
		p.DeleteCallback(p.Game.Code)
		return
	} else if len(p.Game.Players) < 2 {
		p.Game.State = game.STATE_GAME_LOBBY
		p.Game.CurrentRound = nil
		p.sendLobbyState(game.STATE_GAME_LOBBY, false)
		return
	}

	if p.Game.Mod == player.Name {
		p.Game.Mod = p.Game.Players[0].Name
	}

	p.Game.UpdateState()

	if p.Game.CurrentRound.Boss == player.Name {
		p.Game.StartNewRound()
		p.sendPlayersChoosingState()
	} else {
		p.Game.CurrentRound.RemoveAllForPlayer(player)
		p.sendCurrentState(playerName)
	}

	p.Game.UpdateState()
}

func (p *Play) handlePlayerKick(playerName string, message any) {
	if p.Game.Mod != playerName {
		log.Printf("Player %s is not MOD and cannot kick players from game %s!", playerName, p.Game.Code)
		p.sendError(playerName, err.ACTION_FORBIDDEN, game.ACTION_PLAYER_KICK)
		return
	}

	action, ok := message.(communication.PlayerKickAction)
	if !ok {
		p.sendError(playerName, err.BAD_REQUEST, game.ACTION_PLAYER_KICK)
		return
	}

	player := p.Game.FindPlayer(action.PlayerName)
	if player == nil {
		p.sendError(playerName, err.PLAYER_NOT_FOUND, action.PlayerName)
		return
	}

	p.handlePlayerLeave(player.Name)
}

func (p *Play) sendCurrentState(playerName string) {
	log.Printf("Sending current state %s for game %s!", p.Game.State, p.Game.Code)

	switch p.Game.State {
	case game.STATE_GAME_LOBBY:
		p.sendLobbyState(game.STATE_GAME_LOBBY, false)
	case game.STATE_GAME_READY:
		p.sendLobbyState(game.STATE_GAME_LOBBY, true)
	case game.STATE_PLAYERS_CHOOSING:
		p.sendPlayersChoosingState()
	case game.STATE_BOSS_CHOOSING:
		p.sendBossChoosingState()
	case game.STATE_ROUND_FINISHED:
		p.sendBossHasChosenState(p.Game.CurrentRound.Winner, p.Game.CurrentRound.PlayedCards[p.Game.CurrentRound.Winner])
	case game.STATE_GAME_FINISHED:
		p.sendLobbyState(game.STATE_GAME_FINISHED, true)
	default:
		p.sendError(playerName, err.INVALID_STATE, p.Game.State)
	}
}

func (p *Play) sendError(playerName string, label string, payload any) {
	state := &communication.ApplicationError{
		Label:   label,
		Payload: payload,
	}

	p.senders[playerName].Send("game.error", state)
	log.Printf("Sent error '%s' to player %s for game %s!", label, playerName, p.Game.Code)
}
