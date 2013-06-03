// appends png suffixes to graphite urls  so they render in the irc clients
package graphite

import (
	"fmt"
	"github.com/collinvandyck/gesture/core"
	"log"
	"strings"
)

func Create(bot *core.Gobot, config map[string]interface{}) {
	prefix, found := config["prefix"].(string)
	if !found {
		log.Printf("Can't find graphite prefix!")
		return
	}

	pattern := fmt.Sprintf(`%s(\S+)`, prefix)
	bot.ListenFor(pattern, func(msg core.Message, matches []string) core.Response {
		url := matches[0]
		if !strings.HasSuffix(url, ".png") {
			msg.Ftfy(url + "&lol.png")
		}
		return bot.Stop()
	})
}
