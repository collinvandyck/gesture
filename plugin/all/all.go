// mentions everyone in the group so that they get notified about a message
package all

import (
	"github.com/collinvandyck/gesture/core"
	"sort"
	"strings"
)

func Create(bot *core.Gobot, config map[string]interface{}) {
	bot.ListenFor("!all", func(msg core.Message, matches []string) core.Response {
		names := make([]string, 0)
		for _, name := range msg.Names() {
			if name != msg.User && name != bot.Name {
				names = append(names, name)
			}
		}
		if len(names) > 0 {
			sort.Strings(names)
			msg.Send("cc: " + strings.Join(names, " "))
		}
		return bot.Stop()
	})
}
