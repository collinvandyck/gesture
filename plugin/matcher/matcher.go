// the matcher plugin runs completely from configuration, subsituting regular expressions
// for static strings.
package matcher

import (
	"github.com/collinvandyck/gesture/core"
	"log"
)

func Create(bot *core.Gobot, config map[string]interface{}) {
	// matches is actually a map[string]string
	matches, found := config["matches"]
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
