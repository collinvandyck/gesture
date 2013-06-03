package sensu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/collinvandyck/gesture/core"
	"github.com/collinvandyck/gesture/util"
	"log"
	"net/http"
	"strings"
	"time"
)

type sensuEventData struct {
	Client      string
	Check       string
	Occurrences uint64
	Output      string
	Flapping    bool
	Status      uint8
}

type sensuEvent interface {
	statusAsString()
}

type eventsResponse []sensuEventData

type stashList []string

type postData struct {
	timestamp int64
}

func Create(bot *core.Gobot, config map[string]interface{}) {
	if len(config) == 0 {
		log.Printf("No sensu config found")
		return
	}

	envs := make(map[string]string)
	for e, u := range config["environments"].(map[string]interface{}) {
		envs[e] = fmt.Sprintf("%s", u)
	}

	bot.ListenFor("^sensu (.*)", func(msg core.Message, matches []string) core.Response {
		cmdArgs := strings.Split(matches[1], " ")
		switch cmdArgs[0] {
		case "events":
			if len(cmdArgs) > 1 {
				if err, events := getEvents(envs[cmdArgs[1]]); err != nil {
					return bot.Error(err)
				} else {
					msg.Send(fmt.Sprintf("%s: Total events: %d.", cmdArgs[1], len(events)))

					for _, event := range events {
						msg.SendPriv(fmt.Sprintf("%s: %s", cmdArgs[1], event.toString()))
					}
				}
			} else {
				for env, url := range envs {
					if err, events := getEvents(url); err != nil {
						return bot.Error(err)
					} else {
						msg.Send(fmt.Sprintf("%s: Total events: %d.", env, len(events)))

						for _, event := range events {
							msg.SendPriv(fmt.Sprintf("%s: %s", env, event.toString()))
						}
					}
				}
			}
		case "silence":
			var (
				env    string
				target string
			)

			if len(cmdArgs) > 2 {
				env = envs[cmdArgs[1]]
				target = cmdArgs[2]
			} else {
				env = ""
				target = cmdArgs[1]
			}

			if err := silence(env, target); err != nil {
				return bot.Error(err)
			} else {
				msg.Send(fmt.Sprintf("silenced %s in env: %s", cmdArgs[2], cmdArgs[1]))
			}
		case "silenced":
			if len(cmdArgs) > 1 {
				if err, silenced := getSilenced(envs[cmdArgs[1]]); err != nil {
					return bot.Error(err)
				} else {
					msg.Send(fmt.Sprintf("%s: Total silenced: %d.", cmdArgs[1], len(silenced)))

					for _, s := range silenced {
						msg.SendPriv(fmt.Sprintf("%s: %s", cmdArgs[1], s))
					}
				}
			} else {
				for env, url := range envs {
					if err, silenced := getSilenced(url); err != nil {
						return bot.Error(err)
					} else {
						msg.Send(fmt.Sprintf("%s: Total silenced: %d.", env, len(silenced)))

						for _, s := range silenced {
							msg.SendPriv(fmt.Sprintf("%s: %s", env, s))
						}
					}
				}
			}
		case "unsilence":
			var (
				env    string
				target string
			)

			if len(cmdArgs) > 2 {
				env = envs[cmdArgs[1]]
				target = cmdArgs[2]
			} else {
				env = ""
				target = cmdArgs[1]
			}

			if err := unsilence(env, target); err != nil {
				fmt.Println(err)
			} else {
				msg.Send(fmt.Sprintf("silenced %s in env: %s", cmdArgs[2], cmdArgs[1]))
			}
		}

		return bot.Stop()
	})
}

func getEvents(sensuUrl string) (error, eventsResponse) {
	eventsUrl := fmt.Sprintf("%s/events", sensuUrl)
	var eventsResponse eventsResponse
	err := util.UnmarshalUrl(eventsUrl, &eventsResponse)
	return err, eventsResponse
}

func getStashes(sensuUrl string) (error, stashList) {
	stashesUrl := fmt.Sprintf("%s/stashes", sensuUrl)
	var stashResponse stashList
	err := util.UnmarshalUrl(stashesUrl, &stashResponse)
	return err, stashResponse
}

func getSilenced(sensuUrl string) (error, []string) {
	var silenced []string
	if err, stashes := getStashes(sensuUrl); err != nil {
		return err, silenced
	} else {
		for _, stash := range stashes {
			if strings.HasPrefix(stash, "silence/") {
				silenced = append(silenced, string(stash[8:]))
			}
		}
	}
	return nil, silenced
}

func silence(sensuUrl string, target string) error {
	data := postData{timestamp: time.Now().Unix()}
	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	silenceUrl := fmt.Sprintf("%s/stash/silence/%s", sensuUrl, target)
	_, err = http.Post(silenceUrl, "application/json", bytes.NewBuffer(marshalled))
	return err
}

func unsilence(sensuUrl string, target string) error {
	silenceUrl := fmt.Sprintf("%s/stash/silence/%s", sensuUrl, target)
	req, _ := http.NewRequest("DELETE", silenceUrl, nil)
	_, err := http.DefaultClient.Do(req)
	return err
}

func (event *sensuEventData) statusAsString() string {
	var status string
	switch event.Status {
	case 0:
		status = "OK"
	case 1:
		status = "WARNING"
	case 2:
		status = "CRITICAL"
	case 3:
		status = "UNKNOWN"
	}
	return status
}

func (event *sensuEventData) toString() string {
	return fmt.Sprintf("%s: %s/%s - %s", event.statusAsString(), event.Client, event.Check, event.Output)
}
