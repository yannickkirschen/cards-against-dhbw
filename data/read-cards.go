package data

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/yannickkirschen/cards-against-dhbw/model"
)

func ReadCards(blacks []*model.Card, whites []*model.Card) error {
	card, err := ioutil.ReadFile("cards.json")
	if err != nil {
		log.Println("failed to read json: " + err.Error())
		return err
	}

	var data []*model.Card
	err = json.Unmarshal(card, &data)

	if err != nil {
		log.Fatalln("failed to parse cards: ", err)
		return nil
	}

	for _, s := range data {
		if s.Type == model.BLACK {
			blacks = append(blacks, s)
		} else {
			whites = append(whites, s)
		}
	}

	return nil
}
