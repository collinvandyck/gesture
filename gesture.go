package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gesture/plugin"
	"gesture/plugin/identity"
	"gesture/plugin/gis"
	"gesture/plugin/twitter"
	"gesture/rewrite"
	irc "github.com/fluffle/goirc/client"
	"log"
	"os"
	"strings"
	"io/ioutil"
)

var (
	plugins  []plugin.Plugin
)

// gesture config
type Config struct {
	BotName  string
	Hostname string
	SSL bool
	Channels []string
}

// readsConfig unmarshals the config from a file and returns the struct
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

// a Plugin is something that can respond to messages
func main() {
	if len(os.Args) < 2 {
		log.Println("usage: gesture [conf_file]")
		os.Exit(1)
	}

	config, err := readConfig(os.Args[1])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	plugins = []plugin.Plugin{
		twitter.NewPlugin(), 
		gis.NewPlugin(),
		identity.NewPlugin(config.BotName),
	}

	flag.Parse()
	c := irc.SimpleClient(config.BotName)
	c.SSL = config.SSL
	c.AddHandler(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("Connected to", config.Hostname, "!")
			for _, channel := range config.Channels {
				conn.Join(channel)
			}
		})
	quit := make(chan bool)
	c.AddHandler("JOIN", func(conn *irc.Conn, line *irc.Line) { 
		if line.Nick == config.BotName {
			log.Printf("Joined %+v\n", line.Args)
		}
	})
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { 
		log.Println("Disconnected. Quitting.")
		quit <- true 
	})
	c.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		messageReceived(conn, line)
	})
	if err := c.Connect(config.Hostname); err != nil {
		log.Fatalf("Connection error: %s\n", err)
	}
	// Wait for disconnect
	<-quit
}

// When a message comes in on a channel gesture has joined, this method will be called.
func messageReceived(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) > 1 {
		channel := line.Args[0]
		message := line.Args[1]

		mc := &messageContext{conn, line}

		log.Printf(">> %s (%s): %s\n", line.Nick, channel, message)

		handled := false
		for _, plugin := range plugins {
			success, err := plugin.Call(mc)
			if err != nil {
				log.Print(err)
			}
			if success {
				handled = true
				break
			}
		}
		if !handled {
			// try to expand any links
			for _, token := range rewrite.GetRewrittenLinks(mc.Message()) {
				mc.Reply(token)
			}
		}
	}
}

// when an error occurs, calling this method will send the error back to the irc channel
func sendError(conn *irc.Conn, channel string, nick string, err error) {
	log.Print(err)
	conn.Privmsg(channel, fmt.Sprintf("%s: oops: %v", nick, err))
}

type messageContext struct {
	conn *irc.Conn
	line *irc.Line
}

func (mc *messageContext) Message() string {
	if len(mc.line.Args) > 1 {
		return mc.line.Args[1]
	}
	return ""
}

func (mc *messageContext) Command() string {
	sliced := strings.Split(mc.Message(), " ")
	return sliced[0]
}

func (mc *messageContext) CommandArgs() []string {
	sliced := strings.Split(mc.Message(), " ")
	return sliced[1:]
}

func (mc *messageContext) Reply(message string) {
	channel := mc.line.Args[0]
	mc.conn.Privmsg(channel, fmt.Sprintf("%s: ftfy -> %s", mc.line.Nick, rewrite.Rewrite(message)))
}

func (mc *messageContext) Send(message string) {
	channel := mc.line.Args[0]
	mc.conn.Privmsg(channel, rewrite.Rewrite(message))
}
