// does twitter related things
package twitter

import (
	"encoding/json"
	"fmt"
	"gesture/core"
	"gesture/util"
	"strings"
)

func Create(bot *core.Gobot) {
	bot.ListenFor("^describe (\\w+)", func(msg core.Message, matches []string) error {
		described, err := describe(matches[1])
		if err == nil {
			msg.Send(described)
		}
		return err
	})

	bot.ListenFor("twitter\\.com/(\\w+)/status/(\\d+)", func(msg core.Message, matches []string) error {
		user, tweet, err := getTweet(matches[2])
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
		return err
	})
}

func getTweet(tweetId string) (user string, tweet string, err error) {
	body, err := util.GetUrl("https://api.twitter.com/1/statuses/show/" + tweetId + ".json")
	if err != nil {
		return "", "", err
	}
	var content map[string]interface{}
	if err = json.Unmarshal(body, &content); err != nil {
		return "", "", err
	}

	user = content["user"].(map[string]interface{})["screen_name"].(string)
	tweet = content["text"].(string)

	return user, tweet, nil
}

func describe(user string) (result string, err error) {
	body, err := util.GetUrl("http://api.twitter.com/1/users/lookup.json?include_entities=true&screen_name=" + user)
	if err != nil {
		return "", err
	}
	var jsonResponse []map[string]interface{}
	if err = json.Unmarshal(body, &jsonResponse); err != nil {
		return "", err
	}
	first := jsonResponse[0]
	description := first["description"].(string)
	pic := first["profile_image_url_https"].(string)
	return fmt.Sprintf("\"%s\" %s", description, pic), nil
}
