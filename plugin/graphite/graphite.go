// appends png suffixes to graphite urls  so they render in the irc clients
package graphite

import (
	"gesture/plugin"
	"strings"
)

type Plugin struct {
	graphitePrefix string
}

func NewPlugin(graphitePrefix string) Plugin {
	return Plugin{graphitePrefix}
}

func (me Plugin) Call(mc plugin.MessageContext) (bool, error) {
	for _, token := range strings.Split(mc.Message(), " ") {
		if strings.HasPrefix(token, me.graphitePrefix) {
			if !strings.HasSuffix(token, ".png") {
				mc.Ftfy(token + "&lol.png")
			}
		}
	}
	return false, nil
}
