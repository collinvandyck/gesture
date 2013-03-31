// Talkin'
package core

import (
	"fmt"
	"gesture/rewrite"
	irc "github.com/fluffle/goirc/client"
)

type Message struct {
	conn *irc.Conn
	line *irc.Line
	User string
	Channel string
	Text string
}

func (msg *Message) Send(message string) {
	msg.conn.Privmsg(msg.Channel, rewrite.Rewrite(message))
}

func (msg *Message) Reply(message string) {
	msg.Send(fmt.Sprintf("%s: %s", msg.User, message))
}

func (msg *Message) Ftfy(message string) {
	msg.Send(fmt.Sprintf("%s: ftfy -> %s", msg.User, message))
}
