// DO IT
package main

import (
	"flag"
	"github.com/collinvandyck/gesture/core"
	"github.com/collinvandyck/gesture/plugin/all"
	"github.com/collinvandyck/gesture/plugin/gis"
	"github.com/collinvandyck/gesture/plugin/graphite"
	"github.com/collinvandyck/gesture/plugin/identity"
	"github.com/collinvandyck/gesture/plugin/markov"
	"github.com/collinvandyck/gesture/plugin/matcher"
	"github.com/collinvandyck/gesture/plugin/memegenerator"
	"github.com/collinvandyck/gesture/plugin/sensu"
	"github.com/collinvandyck/gesture/plugin/twitter"
	"github.com/collinvandyck/gesture/plugin/youtube"
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
	rand.Seed(time.Now().UTC().UnixNano())

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
