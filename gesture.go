package main

import (
	"flag"
	"fmt"
	"gesture/plugin"
	"gesture/plugin/gis"
	"gesture/rewrite"
	"gesture/plugin/twitter"
	irc "github.com/fluffle/goirc/client"
	"log"
	"strings"
)

var (
	channels = []string{"#collinjester"}
	plugins  []plugin.Plugin
)

// a Plugin is something that can respond to messages
func main() {
	var twitterPlugin = twitter.NewPlugin()
	var gisPlugin = gis.NewPlugin()
	plugins = []plugin.Plugin{twitterPlugin, gisPlugin}

	flag.Parse()
	c := irc.SimpleClient("gesturebot")
	c.SSL = true
	c.AddHandler(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			for _, channel := range channels {
				conn.Join(channel)
			}
		})
	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { quit <- true })
	c.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		messageReceived(conn, line)
	})
	if err := c.Connect("irc.freenode.net"); err != nil {
		fmt.Printf("Connection error: %s\n", err)
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
			// try to rewrite the line
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
	mc.conn.Privmsg(channel, fmt.Sprintf(rewrite.Rewrite(message)))
}



