// mentions everyone in the group so that they get notified about a message
package all

import (
	"gesture/core"
	"sort"
	"strings"
)

func Create(bot *core.Gobot) {
	bot.ListenFor("!all", func(msg core.Message, matches []string) error {
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
		return nil
	})
}
