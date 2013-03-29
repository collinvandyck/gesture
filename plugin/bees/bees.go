package bees

import (
	"gesture/plugin"
	"strings"
)

type Plugin bool

func New() Plugin {
	return Plugin(true)
}

func (p Plugin) Call(mc plugin.MessageContext) (bool, error) {
	if strings.Contains(mc.Message(), "bees") {
		mc.Send("http://i.imgur.com/qrLEV.gif")
	}
	return false, nil
}
