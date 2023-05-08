package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	socket "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/yannickkirschen/cards-against-dhbw/communication"
	e "github.com/yannickkirschen/cards-against-dhbw/err"
	"github.com/yannickkirschen/cards-against-dhbw/game"
	"github.com/yannickkirschen/cards-against-dhbw/shelf"
)

type socketSender struct {
	f func(title string, message ...any)
}

func (s *socketSender) Send(title string, message ...any) {
	s.f(title, message...)
}

func InitServerSession() *socket.Server {
	server := socket.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{CheckOrigin: func(r *http.Request) bool { return true }},
			&websocket.Transport{CheckOrigin: func(r *http.Request) bool { return true }},
		},
	})

	server.OnConnect("/", onConnect)
	server.OnEvent("/", game.ACTION_GAME_JOIN, onEventJoin)
	server.OnEvent("/", game.ACTION_GAME_START, onEventGameStart)
	server.OnEvent("/", game.ACTION_CARD_CHOSEN, onEventCardChosen)
	server.OnEvent("/", game.ACTION_ROUND_CONTINUE, onEventRoundContinue)
	server.OnEvent("/", game.ACTION_CARD_REMOVE, onEventRemoveCard)
	server.OnEvent("/", game.ACTION_GAME_LEAVE, onEventLeaveGame)
	server.OnEvent("/", game.ACTION_PLAYER_KICK, onKickEvent)
	server.OnEvent("/", "notice", onEventNotice)
	server.OnError("/", onError)
	server.OnDisconnect("/", onDisconnect)

	go server.Serve()
	return server
}

func onConnect(s socket.Conn) error {
	var p communication.JoinRequestAction = communication.JoinRequestAction{GameCode: "unregistered", PlayerName: "unregistered"}
	s.SetContext(p)

	log.Printf("Socket connection created.")
	return nil
}

func onEventJoin(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction
	json.Unmarshal([]byte(msg), &p)
	s.SetContext(p)

	log.Printf("Game %s: received message from player %s to join the game!", p.GameCode, p.PlayerName)

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot join the game!", p.GameCode, p.PlayerName)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot join the game!", p.GameCode, p.PlayerName)
	}

	err = gp.AddSender(p.PlayerName, &socketSender{f: s.Emit})
	if err != nil {
		s.Emit("game.error", &communication.ApplicationError{
			Label:   e.PLAYER_NOT_FOUND,
			Payload: p.PlayerName,
		})
		return "player not found"
	}

	gp.Receive(p.PlayerName, game.ACTION_GAME_JOIN, msg)
	return fmt.Sprintf("Player %s joined the game!", p.PlayerName)
}

func onEventGameStart(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message from player %s to start the game!", p.GameCode, p.PlayerName)

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot start the game!", p.GameCode, p.PlayerName)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot start the game!", p.GameCode, p.PlayerName)
	}

	gp.Receive(p.PlayerName, game.ACTION_GAME_START, msg)
	return fmt.Sprintf("Player %s started the game!", p.PlayerName)
}

func onEventCardChosen(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message that player %s chose a card!", p.GameCode, p.PlayerName)

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot choose a card!", p.GameCode, p.PlayerName)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot choose a card!", p.GameCode, p.PlayerName)
	}

	var action communication.CardChosenAction
	err = json.Unmarshal([]byte(msg), &action)
	if err != nil {
		log.Printf("Game %s (socket): error while parsing card chosen action: %s", p.GameCode, err)
		return fmt.Sprintf("Game %s (socket): error while parsing card chosen action: %s", p.GameCode, err)
	}

	gp.Receive(p.PlayerName, "entity.card.chosen", action)
	return fmt.Sprintf("Player %s chose a card!", p.PlayerName)
}

func onEventRoundContinue(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message that mod %s chose to continue the round!", p.GameCode, p.PlayerName)

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		log.Printf("Game %s (socket): game not found, mod %s cannot choose to continue the round!", p.GameCode, p.PlayerName)
		return fmt.Sprintf("Game %s (socket): game not found, mod %s cannot choose to continue the round!", p.GameCode, p.PlayerName)
	}

	gp.Receive(p.PlayerName, game.ACTION_ROUND_CONTINUE, msg)
	return fmt.Sprintf("MOD %s chose to continue the round!", p.PlayerName)
}

func onEventRemoveCard(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message that %s wants to remove a card!", p.GameCode, p.PlayerName)

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		log.Printf("Game %s (socket): game not found, %s cannot remove a card!", p.GameCode, p.PlayerName)
		return fmt.Sprintf("Game %s (socket): game not found, %s cannot remove a card!", p.GameCode, p.PlayerName)
	}

	var action communication.RemoveCardAction
	err = json.Unmarshal([]byte(msg), &action)
	if err != nil {
		log.Printf("Game %s (socket): error while parsing remove card action: %s", p.GameCode, err)
		return fmt.Sprintf("Game %s (socket): error while parsing remove card action: %s", p.GameCode, err)
	}

	gp.Receive(p.PlayerName, game.ACTION_CARD_REMOVE, action)
	return fmt.Sprintf("Player %s removed a card!", p.PlayerName)
}

func onEventLeaveGame(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message that player %s wants to leave the game!", p.GameCode, p.PlayerName)

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot leave the game!", p.GameCode, p.PlayerName)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot leave the game!", p.GameCode, p.PlayerName)
	}

	gp.Receive(p.PlayerName, game.ACTION_GAME_LEAVE, msg)
	return fmt.Sprintf("Player %s left the game!", p.PlayerName)
}

func onKickEvent(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message that mod %s kicks a player!", p.GameCode, p.PlayerName)

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot kick another player!", p.GameCode, p.PlayerName)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot kick another player!", p.GameCode, p.PlayerName)
	}

	var action communication.PlayerKickAction
	err = json.Unmarshal([]byte(msg), &action)
	if err != nil {
		log.Printf("Game %s (socket): error while parsing player kick action: %s", p.GameCode, err)
		return fmt.Sprintf("Game %s (socket): error while parsing player kick action: %s", p.GameCode, err)
	}

	gp.Receive(p.PlayerName, game.ACTION_PLAYER_KICK, action)
	return fmt.Sprintf("Player %s kicked another!", p.PlayerName)
}

func onEventNotice(s socket.Conn, msg string) {
	fmt.Println("Printing Context", msg)
	fmt.Println("Context: ", s.Context())
	fmt.Println("Context-Type", reflect.TypeOf(s.Context()))
}

func onError(s socket.Conn, e error) {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): error occurred: %s", p.GameCode, e.Error())

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		return
	}
	if e.Error() == "websocket: close 1001 (going away)" || e.Error() == "websocket: close sent" {
		log.Printf("Ignoring error: %s", e.Error())
		return
	}

	gp.Receive(p.PlayerName, game.ACTION_GAME_LEAVE, "")
}

func onDisconnect(s socket.Conn, reason string) {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): connection %s closed: %s", p.GameCode, s.ID(), reason)

	gp, err := shelf.GlobalShelf.Play(p.GameCode)
	if err != nil {
		return
	}

	gp.Receive(p.PlayerName, "player.inactive", "")
}
