// Reads the `dhbw.json` file and stores the configuration in the global variable `DhbwConfig`.
//
// Author: Yannick Kirschen
package config

import (
	"encoding/json"
	"log"
	"os"
)

// Holds the configuration
type Config struct {
	Database Database
	Port     int `json:"port"`
}

// Holds the database configuration
type Database struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Holds the global configuration
var DhbwConfig Config

// Reads the dhbw.json file
func InitConfig() error {
	log.Print("Attempting to read config file...")

	file, err := os.ReadFile("dhbw.json")
	if err != nil {
		log.Printf("Unable to read config file: %s", err.Error())
		return err
	}

	err = json.Unmarshal(file, &DhbwConfig)
	if err != nil {
		log.Printf("Unable to parse config file: %s", err.Error())
		return err
	}

	log.Print("Successfully read config file!")
	return nil
}
