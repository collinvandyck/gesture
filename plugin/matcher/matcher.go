// the matcher plugin runs completely from configuration, subsituting regular expressions
// for static strings.
package matcher

import (
	"gesture/core"
	"log"
)

func Create(bot *core.Gobot) {
	// matches is actually a map[string]string
	matches, found := bot.Config.Plugins["matcher"]["matches"]
	if !found {
		log.Printf("Can't find matcher/matches plugin conf. Plugin will not run.")
		return
	}

	switch matches := matches.(type) {
	case map[string]interface{}:
		for pattern, replacement := range matches {
			switch replacement := replacement.(type) {
			case string:
				bot.ListenFor(pattern, func(msg core.Message, matches []string) core.Response {
					msg.Send(replacement)
					return bot.KeepGoing()
				})
			}
		}
	}

}
