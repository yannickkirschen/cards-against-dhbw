package main

import (
	"fmt"
	"net/http"

	"github.com/yannickkirschen/cards-against-dhbw/config"
	"github.com/yannickkirschen/cards-against-dhbw/db"
	"github.com/yannickkirschen/cards-against-dhbw/game"
	"github.com/yannickkirschen/cards-against-dhbw/server"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	db.ReadCards()
	game.GlobalGameShelf = game.NewGameShelf()

	http.HandleFunc("/v1/new", game.NewGameHandler)
	http.HandleFunc("/v1/join/", game.JoinGameHandler)

	http.Handle("/socket.io/", server.InitServerSession())

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.DhbwConfig.Port), nil))
}
