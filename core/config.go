// Configuration for your GoBot
// "How to wave your hands"
package core

import (
	"encoding/json"
	"os"
)

type Config struct {
	// By default, goirc enables this
	DisableFloodProtection bool
	BotName                string
	Hostname               string
	SSL                    bool
	Channels               []string
	Plugins                map[string]map[string]interface{}
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	dec := json.NewDecoder(file)
	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
