// Your very own Gobot
// "More man than machine"
package core

import (
	"flag"
	"gesture/rewrite"
	irc "github.com/fluffle/goirc/client"
	"log"
	"regexp"
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
func CreateGobot(conf *Config) *Gobot {
	bot := &Gobot{conf.BotName, conf, nil, make(chan bool), nil}

	bot.setupIrc(conf)
	// TODO: Add enabled/disabled plugins to the config
	var enabled, disabled []string
	bot.loadPlugins(enabled, disabled)

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
func (bot *Gobot) ListenFor(pattern string, cb func(Message, []string) error) {
	re := regexp.MustCompile(pattern)
	bot.listeners = append(bot.listeners, listener{re, cb})
}

// TODO:
// - OnEnter/Leave
// - OnTopicChange

// -------------------------------------------------------------------
// GOBOT'S ROOM, KEEP OUT

func (bot *Gobot) setupIrc(conf *Config) {
	flag.Parse()
	bot.client = irc.SimpleClient(conf.BotName)
	bot.client.SSL = conf.SSL

	bot.client.AddHandler(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("Connected to", conf.Hostname, "!")
			for _, channel := range conf.Channels {
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
}

func (bot *Gobot) loadPlugins(enabled, disabled []string) {
	allPlugins := GetAllPlugins()
	log.Printf("Loading all %d available plugins", len(allPlugins))
	// TODO: Actually do something with enabled/disabled plugins
	for _, p := range allPlugins {
		if err := p.Create(bot); err != nil {
			log.Printf("Failed to initialize %s plugin: %s", p.Name(), err)
		} else {
			log.Printf(" --> %s", p.Name())
		}
	}
}

func (bot *Gobot) messageReceived(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) > 1 {
		msg := messageFrom(conn, line)
		log.Printf(">> %s (%s): %s\n", msg.User, msg.Channel, msg.Text)

		matched := false
		var err error = nil
		for _, listener := range bot.listeners {
			matched, err = listener.listen(msg)
			if err != nil {
				log.Print(err)
			}
			if matched {
				break
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

type listener struct {
	re *regexp.Regexp
	cb func(Message, []string) error
}

// Try to match the given message. If it matches, fire the callback and returns
// true. Returns false otherwise.
func (listener *listener) listen(msg Message) (matched bool, err error) {
	if matches := listener.re.FindStringSubmatch(msg.Text); matches != nil {
		matched = true
		err = listener.cb(msg, matches)
	}
	return
}
