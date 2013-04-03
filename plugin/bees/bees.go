package bees

import (
	"gesture/core"
)

func Create(bot *core.Gobot) {
	bot.ListenFor("bee(e*)s", func(msg core.Message, matches []string) error {
		msg.Send("http://i.imgur.com/qrLEV.gif")
		return nil
	})
}
