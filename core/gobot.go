// Your very own Gobot
// "More man than machine"
package core

import (
	"flag"
	"github.com/collinvandyck/gesture/rewrite"
	irc "github.com/fluffle/goirc/client"
	"log"
	"regexp"
)

type Response struct {
	Status Status
	Error  error
}

type Status int

const (
	Stop Status = iota
	KeepGoing
)

type Gobot struct {
	Name      string
	Config    *Config
	client    *irc.Conn
	quitter   chan bool
	listeners []listener
}

// -----------------------------------------------------------------------------
// Tell Gobot how to be a Real Boy

// Create a new Gobot from the given gesture config
func CreateGobot(config *Config) *Gobot {
	bot := &Gobot{config.BotName, config, nil, make(chan bool), nil}

	flag.Parse()
	bot.client = irc.SimpleClient(config.BotName)
	bot.client.SSL = config.SSL
	bot.client.Flood = config.DisableFloodProtection
	bot.client.EnableStateTracking()

	bot.client.AddHandler(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("Connected to", config.Hostname, "!")
			for _, channel := range config.Channels {
				conn.Join(channel)
			}
		})

	bot.client.AddHandler("JOIN", func(conn *irc.Conn, line *irc.Line) {
		if line.Nick == bot.Name {
			log.Printf("Joined %+v\n", line.Args)
		}
	})

	bot.client.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		bot.quitter <- true
	})

	bot.client.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		bot.messageReceived(conn, line)
	})

	return bot
}

// Attempt to connect to IRC!
func (bot *Gobot) Connect(hostname string) (chan bool, error) {
	err := bot.client.Connect(bot.Config.Hostname)
	if err != nil {
		return nil, err
	}
	return bot.quitter, nil
}

// Send a disconnect message to your robot
func (bot *Gobot) Disconnect() {
	bot.quitter <- true
}

// Add a listener that matches incoming messages based on the given regexp.
// Matched messages and any submatches are returned to the callback.
func (bot *Gobot) ListenFor(pattern string, cb func(Message, []string) Response) {
	bot.ListenForExcluding(pattern, cb, nil)
}

// Add a listener that matches incoming messages based on the given regexp.
// Matched messages and any submatches are returned to the callback. Any
// messages on channels listed in excludes will be ignored.
func (bot *Gobot) ListenForExcluding(pattern string, cb func(Message, []string) Response, excludes []string) {
	re := regexp.MustCompile(pattern)
	bot.listeners = append(bot.listeners, newListener(re, cb, excludes))
}

func (msg *Gobot) Stop() Response {
	return Response{Stop, nil}
}

func (msg *Gobot) KeepGoing() Response {
	return Response{KeepGoing, nil}
}

func (msg *Gobot) Error(err error) Response {
	return Response{Stop, err}
}

// TODO:
// - OnEnter/Leave
// - OnTopicChange

// -------------------------------------------------------------------
// GOBOT'S ROOM, KEEP OUT

func (bot *Gobot) messageReceived(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) > 1 {
		msg := messageFrom(conn, line)
		log.Printf(">> %s (%s): %s\n", msg.User, msg.Channel, msg.Text)

		matched := false
		for _, listener := range bot.listeners {
			response := listener.listen(msg)
			if response != nil {
				if response.Error != nil {
					log.Print(response.Error)
					msg.Reply(response.Error.Error())
					matched = true
					break
				}
				if response.Status == Stop {
					matched = true
					break
				}
			}
		}
		if !matched {
			// try to expand any links
			for _, token := range rewrite.GetRewrittenLinks(msg.Text) {
				msg.Ftfy(token)
			}
		}
	}
}

func messageFrom(conn *irc.Conn, line *irc.Line) Message {
	return Message{conn, line, line.Nick, line.Args[0], line.Args[1]}
}

// -------------------------------------------------------------------
// PICK UP THE DAMN PHONE

func asSet(strings []string) map[string]bool {
	if strings == nil {
		return nil
	}

	set := make(map[string]bool)
	for _, item := range strings {
		set[item] = true
	}
	return set
}

type listener struct {
	// Match against the name of each channel
	ignoreChannels map[string]bool
	// Match against each message
	re *regexp.Regexp
	// Called with the valid message and the thing that matched
	cb func(Message, []string) Response
}

func newListener(re *regexp.Regexp, cb func(Message, []string) Response,
	excludes []string) listener {
	return listener{ignoreChannels: asSet(excludes), cb: cb, re: re}
}

// Try to match the given message if it does not come in on an excluded channel
// If the message matches, fire the callback and return the response. Returns
// nil otherwise.
func (listener *listener) listen(msg Message) *Response {
	if listener.ignoreChannels[msg.Channel] {
		return nil
	}

	if matches := listener.re.FindStringSubmatch(msg.Text); matches != nil {
		response := listener.cb(msg, matches)
		return &response
	}

	return nil
}
