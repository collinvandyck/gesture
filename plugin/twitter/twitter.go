// does twitter related things
package twitter

import (
	"encoding/json"
	"fmt"
	"gesture/core"
	"gesture/util"
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
		tweet, err := getTweet(matches[2])
		if err == nil && tweet != "" {
			msg.Send(tweet)
		}
		return err
	})
}

func getTweet(tweetId string) (result string, err error) {
	body, err := util.GetUrl("https://api.twitter.com/1/statuses/show/" + tweetId + ".json")
	if err != nil {
		return "", err
	}
	var content map[string]interface{}
	if err = json.Unmarshal(body, &content); err != nil {
		return "", err
	}
	
	user := content["user"].(map[string]interface{})["screen_name"]
	tweet := content["text"].(string)

	response := fmt.Sprintf("%s: %s", user, tweet)
	return response, nil
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

