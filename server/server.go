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

func logRequest(p model.JoinRequestAction, action string) {
	log.Println("new", action, "from ", p.PlayerID, " in ", p.GameID)

}

func getGamePlayFromGameID(gameID string) *model.Game {
	return nil
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
		fmt.Println("connected:", s.ID())
		var p model.JoinRequestAction = model.JoinRequestAction{GameID: "unregistered", PlayerID: "unregistered"}
		s.SetContext(p)
		return nil
	})

	server.OnEvent("/", "joinRequestAction", func(s socket.Conn, msg string) {
		var p model.JoinRequestAction
		json.Unmarshal([]byte(msg), &p)
		logRequest(p, "game.join")
		s.SetContext(p)

		gp, err := game.GlobalGameShelf.GamePlay(p.GameID)
		if err != nil {
			return
		}

		gp.AddSender(p.PlayerID, s.Emit)
		gp.ReceiveMessage(p.PlayerID, "game.join", msg)
	})

	server.OnEvent("/", "startGameAction", func(s socket.Conn, msg string) {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)

		logRequest(p, "game.start")

		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "game.start", msg)
	})

	server.OnEvent("/", "cardChosenAction", func(s socket.Conn, msg string) {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		logRequest(p, "cardChosenAction")
		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "cardChosenAction", msg)
	})

	server.OnEvent("/", "leaveGameAction", func(s socket.Conn, msg string) {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		logRequest(p, "leaveGameAction")
		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "leaveGameAction", msg)
	})

	server.OnEvent("/", "notice", func(s socket.Conn, msg string) {
		fmt.Println("Printing Context", msg)
		fmt.Println("Context: ", s.Context())
		fmt.Println("Context-Type", reflect.TypeOf(s.Context()))

	})

	server.OnError("/", func(s socket.Conn, e error) {
		log.Println("error occurred:", e.Error())
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		logRequest(p, "errorAction")
		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "leaveGameAction", "")
	})

	server.OnDisconnect("/", func(s socket.Conn, reason string) {
		log.Println("connection ", s.ID(), " closed: ", reason)
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		logRequest(p, "disconnectAction")
		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "leaveGameAction", "")
	})

	go server.Serve()

	return server
}
