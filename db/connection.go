// Package db provides a session to manage the database connection.
//
// Author: Yannick Kirschen
package db

import (
	"github.com/redis/go-redis/v9"
	"github.com/yannickkirschen/cards-against-dhbw/config"
)

// Holds a connection to the redis database.
// Make sure to call `Create()` before using this variable.
var Client *redis.Client

// Creates a new connection to the redis database.
func Create() {
	conf := config.DhbwConfig.Database

	Client = redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
	})
}

// Closes the connection to the redis database.
func Close() error {
	return Client.Close()
}
