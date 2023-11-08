package config

import (
	"encoding/json"
	"io/ioutil"
	"ogusers-bot/pkg/logging"
)

var (
	C T
)

func init() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		logging.Logger.Error().
			Err(err).
			Msg("Failed to read config.json")
		return
	}

	err = json.Unmarshal(file, &C)
	if err != nil {
		logging.Logger.Error().
			Err(err).
			Msg("Failed to parse config.json")
	}

}
