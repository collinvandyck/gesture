// does twitter related things
package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/collinvandyck/gesture/core"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	httpClient = &http.Client{}
)

func Create(bot *core.Gobot, config map[string]interface{}) {
	token, ok := config["token"].(string)
	if !ok {
		log.Println("Could not find token in config. Twitter plugin won't work")
		return
	}

	bot.ListenFor("^describe (\\w+)", func(msg core.Message, matches []string) core.Response {
		described, err := describe(token, matches[1])
		if err != nil {
			return bot.Error(err)
		}
		msg.Send(described)
		return bot.Stop()
	})

	bot.ListenFor("twitter\\.com/(\\w+)/status/(\\d+)", func(msg core.Message, matches []string) core.Response {
		user, tweet, err := getTweet(token, matches[2])
		if err != nil {
			return bot.Error(err)
		}
		if err == nil && tweet != "" {
			// Split multi-line tweets into separate PRIVMSG calls
			fields := strings.FieldsFunc(tweet, func(r rune) bool {
				return r == '\r' || r == '\n'
			})
			for _, field := range fields {
				if field != "" {
					msg.Send(fmt.Sprintf("%s: %s", user, field))
				}
			}
		}
		return bot.KeepGoing()
	})
}

func getTweet(token, tweetId string) (user string, tweet string, err error) {
	var content map[string]interface{}
	bytes, err := getAuthorizedUrl(token, "https://api.twitter.com/1.1/statuses/show/"+tweetId+".json")
	if err != nil {
		return "", "", err
	}

	err = json.Unmarshal(bytes, &content)
	if err != nil {
		return "", "", err
	}
	user = content["user"].(map[string]interface{})["screen_name"].(string)
	tweet = content["text"].(string)
	return user, tweet, nil
}

func describe(token, user string) (result string, err error) {
	url := "https://api.twitter.com/1.1/users/lookup.json?include_entities=true&screen_name=" + user
	bytes, err := getAuthorizedUrl(token, url)
	if err != nil {
		return "", err
	}
	var jsonResponse []map[string]interface{}
	err = json.Unmarshal(bytes, &jsonResponse)
	if err != nil {
		return "", err
	}
	first := jsonResponse[0]
	description, ok := first["description"].(string)
	if !ok {
		return "", errors.New("No description available")
	}
	pic, ok := first["profile_image_url_https"].(string)
	if !ok {
		return description, nil
	}
	return fmt.Sprintf("\"%s\" %s", description, pic), nil
}

func getAuthorizedUrl(token, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Bad response code: %d", resp.StatusCode))
	}
	return body, nil
}
