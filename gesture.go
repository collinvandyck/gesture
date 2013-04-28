// DO IT
package main

import (
	"flag"
	"gesture/core"
	"gesture/plugin/all"
	"gesture/plugin/gis"
	"gesture/plugin/graphite"
	"gesture/plugin/identity"
	"gesture/plugin/markov"
	"gesture/plugin/matcher"
	"gesture/plugin/memegenerator"
	"gesture/plugin/sensu"
	"gesture/plugin/twitter"
	"gesture/plugin/youtube"
	"log"
	"math/rand"
	"time"
)

func loadPlugins(bot *core.Gobot) {
	gis.Create(bot)
	matcher.Create(bot)
	twitter.Create(bot)
	youtube.Create(bot)
	identity.Create(bot)
	memegenerator.Create(bot)
	graphite.Create(bot)
	sensu.Create(bot)
	all.Create(bot)
	markov.Create(bot)
}

func main() {
	rand.Seed( time.Now().UTC().UnixNano())

	// Parse command-line arguments in logging package
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalln("usage: gesture conf_file")
	}

	config, err := core.ReadConfig(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	bot := core.CreateGobot(config)
	loadPlugins(bot)
	quit, err := bot.Connect(config.Hostname)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}

	<-quit
}
