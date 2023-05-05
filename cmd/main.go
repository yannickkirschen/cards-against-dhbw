package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yannickkirschen/cards-against-dhbw/config"
	"github.com/yannickkirschen/cards-against-dhbw/data"
	"github.com/yannickkirschen/cards-against-dhbw/endpoint"
	"github.com/yannickkirschen/cards-against-dhbw/server"
	"github.com/yannickkirschen/cards-against-dhbw/shelf"
)

func main() {
	log.Print("Welcome to Cards Against DHBW! Starting server ...")

	err := config.InitConfig()
	if err != nil {
		log.Panicf("Unable to start the server: could not read config file. Error was: %s", err.Error())
	}

	data.GlobalClient = data.NewClient()
	shelf.GlobalShelf = shelf.New()

	http.HandleFunc("/v1/new", endpoint.NewGameHandler)
	http.HandleFunc("/v1/join/", endpoint.JoinGameHandler)

	http.Handle("/socket.io/", server.InitServerSession())

	port := config.DhbwConfig.Port
	log.Printf("Server running on port %d", port)
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.DhbwConfig.Port), nil))
}
