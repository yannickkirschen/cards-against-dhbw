package main

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

func InitServerSession() {

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
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "joinRequestAction", func(s socket.Conn, msg string) {
		var p model.JoinRequestAction
		json.Unmarshal([]byte(msg), &p)
		logRequest(p, "joinRequestAction")
		s.SetContext(p)

		getGamePlayFromGameID(p.GameID).SetResponder(p.PlayerID, s.Emit)

		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "joinRequestAction", msg)
	})

	server.OnEvent("/", "startGameAction", func(s socket.Conn, msg string) {
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)

		logRequest(p, "startGameAction")

		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "startGameAction", msg)
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
		s.Emit("reply", "have "+msg)
	})

	server.OnError("/", func(s socket.Conn, e error) {
		log.Fatal("error occurred:", e)
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		logRequest(p, "errorAction")
		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "leaveGameAction", "")
	})

	server.OnDisconnect("/", func(s socket.Conn, reason string) {
		log.Println("connection closed: ", reason)
		var p model.JoinRequestAction = s.Context().(model.JoinRequestAction)
		logRequest(p, "disconnectAction")
		getGamePlayFromGameID(p.GameID).RecvMessage(p.PlayerID, "leaveGameAction", "")
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
