// does twitter related things
package twitter

import (
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
	var content map[string]interface{}
	if err := util.UnmarshalUrl("https://api.twitter.com/1/statuses/show/"+tweetId+".json", &content); err != nil {
		return "", "", err
	}
	user = content["user"].(map[string]interface{})["screen_name"].(string)
	tweet = content["text"].(string)
	return user, tweet, nil
}

func describe(user string) (result string, err error) {
	url := "http://api.twitter.com/1/users/lookup.json?include_entities=true&screen_name=" + user
	var jsonResponse []map[string]interface{}
	if err := util.UnmarshalUrl(url, &jsonResponse); err != nil {
		return "", err
	}
	first := jsonResponse[0]
	description := first["description"].(string)
	pic := first["profile_image_url_https"].(string)
	return fmt.Sprintf("\"%s\" %s", description, pic), nil
}
