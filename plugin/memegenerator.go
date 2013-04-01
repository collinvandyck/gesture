package plugin 

import (
	"encoding/json"
	"errors"
	"gesture/core"
	"gesture/util"
	neturl "net/url"
	"strconv"
)

type Memegenerator struct{}

func init() {
	core.Register(Memegenerator{})
}

func (mg Memegenerator) Name() string {
	return "memegenerator"
}

func (mg Memegenerator) Create(bot *core.Gobot) error {
	username, password, err := loadCredentials(bot.Config.Plugins["memegenerator"])
	if err != nil {
		return err
	}

	fry := memeGen{username, password, fryGenerator, fryImage}

	bot.ListenFor(`(?i)(not sure|unsure) if (.*) or (.*)`, func(msg core.Message, matches []string) error {
		result, err := fry.generate(matches[1]+" if "+matches[2], " or "+matches[3])
		if err == nil && result != "" {
			msg.Reply(result)
		}
		return err
	})

	return nil
}

func loadCredentials(config map[string]interface{}) (string, string, error) {
	user, userOk := config["username"].(string)
	if !userOk {
		return "", "", errors.New("Couldn't find memegenerator username!")
	}
	pass, passOk := config["password"].(string)
	if !passOk {
		return "", "", errors.New("Couldn't find memegenerator password!")
	}
	return user, pass, nil
}

// -----------------------------------------------------------------------------
// Mememememememeeeeeeees

var (
	fryGenerator = 305
	fryImage     = 84688
)

type memeGen struct {
	user      string
	pass      string
	generator int
	image     int
}

func (mg memeGen) generate(firstMsg string, secondMsg string) (string, error) {
	url := "http://version1.api.memegenerator.net/Instance_Create"
	url = url + "?username=" + mg.user
	url = url + "&password=" + mg.pass
	url = url + "&languageCode=en"
	url = url + "&generatorID=" + strconv.Itoa(mg.generator)
	url = url + "&imageID=" + strconv.Itoa(mg.image)
	url = url + "&text0=" + neturl.QueryEscape(firstMsg)
	url = url + "&text1=" + neturl.QueryEscape(secondMsg)

	body, err := util.GetUrl(url)
	if err != nil {
		return "", err
	}

	var decoded map[string]interface{}
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		return "", err
	}
	if result := decoded["result"]; result != nil {
		switch result := result.(type) {
		case map[string]interface{}:
			switch image := result["instanceImageUrl"].(type) {
			case string:
				return image, nil
			}
		}
	}
	return "", nil
}
