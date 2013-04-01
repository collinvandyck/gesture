package plugin 

import (
	"gesture/core"
)

func init() {
	core.Register(Bees{})
}

type Bees struct{}

func (bees Bees) Name() string {
	return "bees"
}

func (bees Bees) Create(bot *core.Gobot) error {
	bot.ListenFor("bee(e*)s", func(msg core.Message, matches []string) error {
		msg.Send("http://i.imgur.com/qrLEV.gif")
		return nil
	})

	return nil
}
