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
	"github.com/yannickkirschen/cards-against-dhbw/game"
	"github.com/yannickkirschen/cards-against-dhbw/model"
)

func allowOriginFunc(r *http.Request) bool {
	return true
}

func InitServerSession() *socket.Server {

	server := socket.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	server.OnConnect("/", func(s socket.Conn) error {
		var p model.JoinRequestAction = model.JoinRequestAction{GameID: "unregistered", PlayerID: "unregistered"}
		s.SetContext(p)

		log.Printf("Game %s (socket): player %s connected!", p.GameID, p.PlayerID)
		return nil
	})

	server.OnEvent("/", "join", func(s socket.Conn, msg string) string {
		var p model.JoinRequestAction
		json.Unmarshal([]byte(msg), &p)
		s.SetContext(p)

		log.Printf("Game %s: received message from player %s to join the game!", p.GameID, p.PlayerID)

		gp, err := game.GlobalGameShelf.GamePlay(p.GameID)
		if err != nil {
			log.Printf("Game %s (socket): game not found, player %s cannot join the game!", p.GameID, p.PlayerID)
			return fmt.Sprintf("Game %s (socket): game not found, player %s cannot join the game!", p.GameID, p.PlayerID)
		}

		gp.AddSender(p.PlayerID, s.Emit)
		gp.ReceiveMessage(p.PlayerID, "game.join", msg)
		return fmt.Sprintf("Player %s joined the game!", p.PlayerID)
	})

	server.OnEvent("/", "game.start", func(s socket.Conn, msg string) string {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		log.Printf("Game %s (socket): received message from player %s to start the game!", p.GameID, p.PlayerID)

		gp, err := game.GlobalGameShelf.GamePlay(p.GameID)
		if err != nil {
			log.Printf("Game %s (socket): game not found, player %s cannot start the game!", p.GameID, p.PlayerID)
			return fmt.Sprintf("Game %s (socket): game not found, player %s cannot start the game!", p.GameID, p.PlayerID)
		}

		gp.ReceiveMessage(p.PlayerID, "game.start", msg)
		return fmt.Sprintf("Player %s started the game!", p.PlayerID)
	})

	server.OnEvent("/", "cardChosenAction", func(s socket.Conn, msg string) string {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		log.Printf("Game %s (socket): received message that player %s chose a card!", p.GameID, p.PlayerID)

		gp, err := game.GlobalGameShelf.GamePlay(p.GameID)
		if err != nil {
			log.Printf("Game %s (socket): game not found, player %s cannot choose a card!", p.GameID, p.PlayerID)
			return fmt.Sprintf("Game %s (socket): game not found, player %s cannot choose a card!", p.GameID, p.PlayerID)
		}

		var action model.CardChosenAction
		err = json.Unmarshal([]byte(msg), &action)
		if err != nil {
			log.Printf("Game %s (socket): error while parsing card chosen action: %s", p.GameID, err)
			return fmt.Sprintf("Game %s (socket): error while parsing card chosen action: %s", p.GameID, err)
		}

		gp.ReceiveMessage(p.PlayerID, "entity.card.chosen", action)
		return fmt.Sprintf("Player %s chose a card!", p.PlayerID)
	})

	server.OnEvent("/", "mod.round.continue", func(s socket.Conn, msg string) string {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		log.Printf("Game %s (socket): received message that mod %s chose to continue the round!", p.GameID, p.PlayerID)

		gp, err := game.GlobalGameShelf.GamePlay(p.GameID)
		if err != nil {
			log.Printf("Game %s (socket): game not found, mod %s cannot choose to continue the round!", p.GameID, p.PlayerID)
			return fmt.Sprintf("Game %s (socket): game not found, mod %s cannot choose to continue the round!", p.GameID, p.PlayerID)
		}

		gp.ReceiveMessage(p.PlayerID, "mod.round.continue", msg)
		return fmt.Sprintf("MOD %s chose to continue the round!", p.PlayerID)
	})

	server.OnEvent("/", "leaveGameAction", func(s socket.Conn, msg string) string {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		log.Printf("Game %s (socket): received message that player %s wants to leave the game!", p.GameID, p.PlayerID)

		gp, err := game.GlobalGameShelf.GamePlay(p.GameID)
		if err != nil {
			log.Printf("Game %s (socket): game not found, player %s cannot leave the game!", p.GameID, p.PlayerID)
			return fmt.Sprintf("Game %s (socket): game not found, player %s cannot leave the game!", p.GameID, p.PlayerID)
		}

		gp.ReceiveMessage(p.PlayerID, "game.leave", msg)
		return fmt.Sprintf("Player %s left the game!", p.PlayerID)
	})

	server.OnEvent("/", "notice", func(s socket.Conn, msg string) {
		fmt.Println("Printing Context", msg)
		fmt.Println("Context: ", s.Context())
		fmt.Println("Context-Type", reflect.TypeOf(s.Context()))

	})

	server.OnError("/", func(s socket.Conn, e error) {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		log.Printf("Game %s (socket): error occurred: %s", p.GameID, e.Error())

		gp, err := game.GlobalGameShelf.GamePlay(p.GameID)
		if err != nil {
			return
		}

		gp.ReceiveMessage(p.PlayerID, "game.leave", "")
	})

	server.OnDisconnect("/", func(s socket.Conn, reason string) {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		log.Printf("Game %s (socket): connection %s closed: %s", p.GameID, s.ID(), reason)

		gp, err := game.GlobalGameShelf.GamePlay(p.GameID)
		if err != nil {
			return
		}

		gp.ReceiveMessage(p.PlayerID, "game.leave", "")
	})

	go server.Serve()
	return server
}
