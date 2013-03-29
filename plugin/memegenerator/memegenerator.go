package memegenerator

import (
	"encoding/json"
	"gesture/plugin"
	"gesture/util"
	neturl "net/url"
	"regexp"
	"strconv"
)

type mg struct {
	username string
	password string
}

var (
	notSureIf = regexp.MustCompile(`(?i)(not sure|unsure) if (.*) or (.*)`)
)

func New(username, password string) plugin.Plugin {
	return mg{username, password}
}

func (mg mg) Call(mc plugin.MessageContext) (bool, error) {
	if match := notSureIf.FindStringSubmatch(mc.Message()); match != nil {
		return generate(mg, mc, 305, 84688, match[1] + " if "+match[2], "or "+match[3])
	}
	return false, nil
}

func generate(mg mg, mc plugin.MessageContext, generatorId int, imageId int, msg1 string, msg2 string) (bool, error) {
	url := "http://version1.api.memegenerator.net/Instance_Create"
	url = url + "?username=" + mg.username
	url = url + "&password=" + mg.password
	url = url + "&languageCode=en"
	url = url + "&generatorID=" + strconv.Itoa(generatorId)
	url = url + "&imageID=" + strconv.Itoa(imageId)
	url = url + "&text0=" + neturl.QueryEscape(msg1)
	url = url + "&text1=" + neturl.QueryEscape(msg2)
	body, err := util.GetUrl(url)
	if err != nil {
		return false, err
	}
	var decoded map[string]interface{}
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		return false, err
	}
	if result := decoded["result"]; result != nil {
		switch result := result.(type) {
		case map[string]interface{}:
			switch image := result["instanceImageUrl"].(type) {
			case string:
				mc.Reply(image)
				return true, nil
			}
		}
	}
	return false, nil
}
