// replies when someone mentions the bot's name
package identity

import (
	"fmt"
	"gesture/core"
)

func Create(bot *core.Gobot) {
	name := bot.Name

	bot.ListenFor(fmt.Sprintf("kill %s", name), func(msg core.Message, matches []string) error {
		msg.Reply("EAT SHIT")
		return nil
	})

	bot.ListenFor(fmt.Sprintf("(hey|h(a?)i|hello) %s", name), func(msg core.Message, matches []string) error {
		msg.Send(fmt.Sprintf("why, hello there %s", msg.User))
		return nil
	})
}
