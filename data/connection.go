// Package db provides a session to manage the database connection.
//
// Author: Yannick Kirschen
package data

import (
	"github.com/redis/go-redis/v9"
	"github.com/yannickkirschen/cards-against-dhbw/config"
)

// Creates a new connection to the redis database.
func NewClient() *redis.Client {
	conf := config.DhbwConfig.Database

	return redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
	})
}

// Closes the connection to the redis database.
func Close(client *redis.Client) error {
	return client.Close()
}
