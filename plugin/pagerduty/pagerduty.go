package pagerduty

import (
	"fmt"
	"gesture/core"
	"github.com/ohlol/pagerduty"
	"log"
)

func Create(bot *core.Gobot) {
	config, found := bot.Config.Plugins["pagerduty"]
	if !found {
		log.Printf("No pagerduty config found")
		return
	}

	account := pagerduty.SetupAccount(config["subdomain"].(string), config["apiKey"].(string))

	bot.ListenFor("^pd (.*)", func(msg core.Message, matches []string) core.Response {
		switch matches[1] {
		case "incidents":
			params := map[string]string {
				"status": "acknowledged,triggered",
			}

			incidents, err := account.Incidents(params)
			if err != nil {
				return bot.Error(err)
			}

			msg.Send(fmt.Sprintf("There are currently %d OPEN (ack,unack) incidents.", len(incidents)))
		}

		return bot.Stop()
	})
}