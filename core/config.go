// Configuration for your GoBot 
// "How to wave your hands"
package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	BotName  string
	Hostname string
	SSL      bool
	Channels []string
	Plugins  map[string]map[string]interface{}
}

func readConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
