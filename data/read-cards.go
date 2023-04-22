package data

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/yannickkirschen/cards-against-dhbw/model"
)

func ReadCards() ([]*model.Card, []*model.Card, error) {
	log.Print("Attempting to read cards file...")

	card, err := ioutil.ReadFile("cards.json")
	if err != nil {
		log.Print("Unable to read cards file: ", err.Error())
		return nil, nil, err
	}

	var data []*model.Card
	err = json.Unmarshal(card, &data)

	if err != nil {
		log.Print("Unable to parse cards file: ", err.Error())
		return nil, nil, err
	}

	blacks := make([]*model.Card, 0)
	whites := make([]*model.Card, 0)
	for _, s := range data {
		if s.Type == model.BLACK {
			blacks = append(blacks, s)
		} else {
			whites = append(whites, s)
		}
	}

	log.Printf("Read %d black cards and %d white cards from cards file.", len(blacks), len(whites))
	return blacks, whites, nil
}
