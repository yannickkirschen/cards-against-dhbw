package data

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/yannickkirschen/cards-against-dhbw/game"
)

// Holds a connection to the redis database.
// Make sure to call `Create()` before using this variable.
var GlobalClient *redis.Client

var Ctx = context.Background()

func Update(game *game.Game) error {
	b, err := json.Marshal(game)
	if err != nil {
		return err
	}

	if GlobalClient.Set(Ctx, game.Code, b, 0).Err() != nil {
		return err
	}

	return nil
}

func FindGame(gameCode string) (*game.Game, error) {
	res, err := GlobalClient.Get(Ctx, gameCode).Result()
	if err == redis.Nil {
		return nil, err
	}

	var game *game.Game

	err = json.Unmarshal([]byte(res), &game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func Delete(gameCode string) error {
	rc, err := GlobalClient.Del(Ctx, gameCode).Result()
	if err != nil {
		return err
	} else if rc == 0 {
		return errors.New("game does not exist in database")
	}

	return nil
}
