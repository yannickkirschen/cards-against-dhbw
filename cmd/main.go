package main

import (
	"io"
	"net/http"

	"github.com/yannickkirschen/cards-against-dhbw/config"
	"github.com/yannickkirschen/cards-against-dhbw/model"
)

var GlobalWhiteCards []model.Card
var GlobalBlackCards []model.Card

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	//http.HandleFunc("/v1/hello", getHello)
	InitServerSession()

}

func getHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, "hello, world\n")
}
