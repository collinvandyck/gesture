// appends png suffixes to graphite urls  so they render in the irc clients
package graphite

import (
	"fmt"
	"gesture/core"
	"log"
	"strings"
)

func Create(bot *core.Gobot) {
	prefix, found := bot.Config.Plugins["graphite"]["prefix"].(string)
	if !found {
		log.Printf("Can't find graphite prefix!")
		return
	}

	pattern := fmt.Sprintf(`%s(\S+)`, prefix)
	bot.ListenFor(pattern, func(msg core.Message, matches []string) error {
		url := matches[0]
		if !strings.HasSuffix(url, ".png") {
			msg.Ftfy(url + "&lol.png")
		}
		return nil
	})
}
