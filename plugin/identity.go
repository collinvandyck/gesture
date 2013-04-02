// replies when someone mentions the bot's name
package plugin 

import (
	"fmt"
	"gesture/core"
)


func init() {
	core.Register(Identity{})
}

type Identity struct{}

func (id Identity) Name() string {
	return "identity"
}

func (id Identity) Create(bot *core.Gobot) error {
	name := bot.Name

	bot.ListenFor(fmt.Sprintf("kill %s", name), func(msg core.Message, matches []string) error {
		msg.Reply("EAT SHIT")
		return nil
	})

	bot.ListenFor(fmt.Sprintf("(hey|h(a?)i|hello) %s", name), func(msg core.Message, matches []string) error {
		msg.Send(fmt.Sprintf("why, hello there %s", msg.User))
		return nil
	})

	return nil
}
