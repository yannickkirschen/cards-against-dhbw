package data

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/yannickkirschen/cards-against-dhbw/card"
)

func ReadCards() ([]*card.Card, []*card.Card, error) {
	log.Print("Attempting to read cards file...")

	c, err := ioutil.ReadFile("cards.json")
	if err != nil {
		log.Print("Unable to read cards file: ", err.Error())
		return nil, nil, err
	}

	var data []*card.Card
	err = json.Unmarshal(c, &data)

	if err != nil {
		log.Print("Unable to parse cards file: ", err.Error())
		return nil, nil, err
	}

	blacks := make([]*card.Card, 0)
	whites := make([]*card.Card, 0)
	for _, s := range data {
		if s.Type == card.BLACK {
			blacks = append(blacks, s)
		} else {
			whites = append(whites, s)
		}
	}

	log.Printf("Read %d black cards and %d white cards from cards file.", len(blacks), len(whites))
	return blacks, whites, nil
}
