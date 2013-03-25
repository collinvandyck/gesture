package main

import (
	"flag"
	"fmt"
	"gesture/gis"
	"gesture/rewrite"
	"gesture/twitter"
	irc "github.com/fluffle/goirc/client"
	"log"
	"strings"
)

var (
	channels = []string{"#collinjester"}
)

// when an error occurs, calling this method will send the error back to the irc channel
func sendError(conn *irc.Conn, channel string, nick string, err error) {
	log.Print(err)
	conn.Privmsg(channel, fmt.Sprintf("%s: oops: %v", nick, err))
}

// When a message comes in on a channel gesture has joined, this method will be called.
func messageReceived(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) > 1 {
		channel := line.Args[0]
		message := line.Args[1]
		messageSliced := strings.Split(message, " ")
		command := messageSliced[0]
		commandArgs := messageSliced[1:]

		log.Printf(">> %s (%s): %s\n", line.Nick, channel, message)

		switch {
		case command == "gis":
			if len(commandArgs) > 0 {
				link, err := gis.Search(strings.Join(commandArgs, " "))
				if err != nil {
					sendError(conn, channel, line.Nick, err)
				} else {
					conn.Privmsg(channel, fmt.Sprintf("%s: %s", line.Nick, link))
				}
			}
		case command == "echo":
			conn.Privmsg(channel, fmt.Sprintf("%s: %s", line.Nick, rewrite.Rewrite(message)))
		case command == "describe":
			if len(commandArgs) > 0 {
				described, err := twitter.Describe(commandArgs[0])
				if err != nil {
					sendError(conn, channel, line.Nick, err)
				} else {
					conn.Privmsg(channel, fmt.Sprintf("%s: %s", line.Nick, described))
				}
			}
		default:
			// find any shortened links and output the expanded versions
			for _, link := range rewrite.GetRewrittenLinks(message) {
				response := line.Nick + ": " + link
				conn.Privmsg(channel, response)
			}
		}
	}
}

func main() {
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
