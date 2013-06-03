package pagerduty

import (
	"fmt"
	"github.com/collinvandyck/gesture/core"
	"github.com/ohlol/pagerduty"
	"log"
)

func Create(bot *core.Gobot, config map[string]interface{}) {
	if len(config) == 0 {
		log.Printf("No pagerduty config found")
		return
	}

	account := pagerduty.SetupAccount(config["subdomain"].(string), config["apiKey"].(string))

	bot.ListenFor("^pd (.*)", func(msg core.Message, matches []string) core.Response {
		switch matches[1] {
		case "incidents":
			params := map[string]string{
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
