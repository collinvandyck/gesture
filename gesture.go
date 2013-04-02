// DO IT
package main

import (
	"encoding/json"
	"gesture/core"
	"gesture/plugin"
	"io/ioutil"
	"log"
	"os"
)

func readConfig(filename string) (*core.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config core.Config
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil

}

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: gesture conf_file")
		os.Exit(1)
	}

	plugin.Init()
	config, err := readConfig(os.Args[1])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	bot := core.CreateGobot(config)
	quit, err := bot.Connect(config.Hostname)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}

	<-quit
}
