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
	server.OnEvent("/", "join", onEventJoin)
	server.OnEvent("/", "game.start", onEventGameStart)
	server.OnEvent("/", "cardChosenAction", onEventCardChosen)
	server.OnEvent("/", "mod.round.continue", onEventRoundContinue)
	server.OnEvent("/", "leaveGameAction", onEventLeaveGame)
	server.OnEvent("/", "notice", onEventNotice)
	server.OnError("/", onError)
	server.OnDisconnect("/", onDisconnect)

	go server.Serve()
	return server
}

func onConnect(s socket.Conn) error {
	var p communication.JoinRequestAction = communication.JoinRequestAction{GameID: "unregistered", PlayerID: "unregistered"}
	s.SetContext(p)

	log.Printf("Game %s (socket): player %s connected!", p.GameID, p.PlayerID)
	return nil
}

func onEventJoin(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction
	json.Unmarshal([]byte(msg), &p)
	s.SetContext(p)

	log.Printf("Game %s: received message from player %s to join the game!", p.GameID, p.PlayerID)

	gp, err := shelf.GlobalShelf.Play(p.GameID)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot join the game!", p.GameID, p.PlayerID)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot join the game!", p.GameID, p.PlayerID)
	}

	gp.AddSender(p.PlayerID, &socketSender{f: s.Emit})
	gp.Receive(p.PlayerID, "game.join", msg)
	return fmt.Sprintf("Player %s joined the game!", p.PlayerID)
}

func onEventGameStart(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message from player %s to start the game!", p.GameID, p.PlayerID)

	gp, err := shelf.GlobalShelf.Play(p.GameID)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot start the game!", p.GameID, p.PlayerID)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot start the game!", p.GameID, p.PlayerID)
	}

	gp.Receive(p.PlayerID, "game.start", msg)
	return fmt.Sprintf("Player %s started the game!", p.PlayerID)
}

func onEventCardChosen(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message that player %s chose a card!", p.GameID, p.PlayerID)

	gp, err := shelf.GlobalShelf.Play(p.GameID)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot choose a card!", p.GameID, p.PlayerID)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot choose a card!", p.GameID, p.PlayerID)
	}

	var action communication.CardChosenAction
	err = json.Unmarshal([]byte(msg), &action)
	if err != nil {
		log.Printf("Game %s (socket): error while parsing card chosen action: %s", p.GameID, err)
		return fmt.Sprintf("Game %s (socket): error while parsing card chosen action: %s", p.GameID, err)
	}

	gp.Receive(p.PlayerID, "entity.card.chosen", action)
	return fmt.Sprintf("Player %s chose a card!", p.PlayerID)
}

func onEventRoundContinue(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message that mod %s chose to continue the round!", p.GameID, p.PlayerID)

	gp, err := shelf.GlobalShelf.Play(p.GameID)
	if err != nil {
		log.Printf("Game %s (socket): game not found, mod %s cannot choose to continue the round!", p.GameID, p.PlayerID)
		return fmt.Sprintf("Game %s (socket): game not found, mod %s cannot choose to continue the round!", p.GameID, p.PlayerID)
	}

	gp.Receive(p.PlayerID, "mod.round.continue", msg)
	return fmt.Sprintf("MOD %s chose to continue the round!", p.PlayerID)
}

func onEventLeaveGame(s socket.Conn, msg string) string {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): received message that player %s wants to leave the game!", p.GameID, p.PlayerID)

	gp, err := shelf.GlobalShelf.Play(p.GameID)
	if err != nil {
		log.Printf("Game %s (socket): game not found, player %s cannot leave the game!", p.GameID, p.PlayerID)
		return fmt.Sprintf("Game %s (socket): game not found, player %s cannot leave the game!", p.GameID, p.PlayerID)
	}

	gp.Receive(p.PlayerID, "game.leave", msg)
	return fmt.Sprintf("Player %s left the game!", p.PlayerID)
}

func onEventNotice(s socket.Conn, msg string) {
	fmt.Println("Printing Context", msg)
	fmt.Println("Context: ", s.Context())
	fmt.Println("Context-Type", reflect.TypeOf(s.Context()))
}

func onError(s socket.Conn, e error) {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): error occurred: %s", p.GameID, e.Error())

	gp, err := shelf.GlobalShelf.Play(p.GameID)
	if err != nil {
		return
	}

	gp.Receive(p.PlayerID, "game.leave", "")
}

func onDisconnect(s socket.Conn, reason string) {
	var p communication.JoinRequestAction = s.Context().(communication.JoinRequestAction)
	log.Printf("Game %s (socket): connection %s closed: %s", p.GameID, s.ID(), reason)

	gp, err := shelf.GlobalShelf.Play(p.GameID)
	if err != nil {
		return
	}

	gp.Receive(p.PlayerID, "game.leave", "")
}
