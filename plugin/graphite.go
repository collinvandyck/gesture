// appends png suffixes to graphite urls  so they render in the irc clients
package plugin 

import (
	"errors"
	"fmt"
	"gesture/core"
	"strings"
)

func init() {
	core.Register(Graphite{})
}

type Graphite struct{}

func (g Graphite) Name() string {
	return "graphite"
}

func (g Graphite) Create(bot *core.Gobot) error {
	prefix, found := bot.Config.Plugins["graphite"]["prefix"].(string)
	if !found {
		return errors.New("Can't find graphite prefix!")
	}

	pattern := fmt.Sprintf(`%s(\S+)`, prefix)
	bot.ListenFor(pattern, func(msg core.Message, matches []string) error {
		url := matches[0]
		if !strings.HasSuffix(url, ".png") {
			msg.Ftfy(url + "&lol.png")
		}
		return nil
	})

	return nil
}
