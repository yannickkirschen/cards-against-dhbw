// Reads the `dhbw.json` file and stores the configuration in the global variable `DhbwConfig`.
//
// Author: Yannick Kirschen
package config

import (
	"encoding/json"
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
	file, fileErr := os.ReadFile("dhbw.json")
	if fileErr != nil {
		return fileErr
	}
	jsonErr := json.Unmarshal(file, &DhbwConfig)
	if jsonErr != nil {
		return jsonErr
	}
	return nil
}
