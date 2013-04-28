// replies when someone mentions the bot's name
package identity

import (
	"fmt"
	"gesture/core"
)

func Create(bot *core.Gobot) {
	name := bot.Name

	bot.ListenFor(fmt.Sprintf("(?i)kill %s", name), func(msg core.Message, matches []string) core.Response {
		msg.Reply("EAT SHIT")
		return bot.Stop()
	})

	bot.ListenFor(fmt.Sprintf("(?i)(hey|h(a?)i|hello) %s", name), func(msg core.Message, matches []string) core.Response {
		msg.Send(fmt.Sprintf("why, hello there %s", msg.User))
		return bot.Stop()
	})
}
