package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yannickkirschen/cards-against-dhbw/config"
	"github.com/yannickkirschen/cards-against-dhbw/game"
	"github.com/yannickkirschen/cards-against-dhbw/server"
)

func main() {
	log.Print("Welcome to Cards Against DHBW! Starting server ...")

	err := config.InitConfig()
	if err != nil {
		log.Panicf("Unable to start the server: could not read config file. Error was: %s", err.Error())
	}

	game.GlobalGameShelf = game.NewGameShelf()

	http.HandleFunc("/v1/new", game.NewGameHandler)
	http.HandleFunc("/v1/join/", game.JoinGameHandler)

	http.Handle("/socket.io/", server.InitServerSession())

	port := config.DhbwConfig.Port
	log.Printf("Server running on port %d", port)
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.DhbwConfig.Port), nil))
}
