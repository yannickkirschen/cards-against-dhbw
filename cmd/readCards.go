package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/yannickkirschen/cards-against-dhbw/model"
)

func readCards() {
	log.Println("----starting to read cards from json-file----")
	cardData, readingError := ioutil.ReadFile("db/cards.json")
	if readingError != nil {
		log.Println("failed to read json: " + readingError.Error())
	}
	var data []*model.Card
	err := json.Unmarshal(cardData, &data)

	if err != nil {
		log.Fatalln("failed to unmarshal cards: ", err)
	}
	for _, s := range data {
		if s.Type == model.BLACK {
			GlobalBlackCards = append(GlobalBlackCards, s)
		} else {
			GlobalWhiteCards = append(GlobalWhiteCards, s)
		}
	}

	/*log.Println("White Cards: ")
	for _, card := range GlobalWhiteCards {
		log.Println("ID:", card.ID)
		log.Println("Text:", card.Text)
		log.Println("------------------------")
	}
	log.Println("Black Cards: ")
	for _, card := range GlobalBlackCards {
		log.Println("ID:", card.ID)
		log.Println("Text:", card.Text)
		log.Println("------------------------")
	}
	*/
	log.Println("----finished reading cards from json-file----")

}
