// Google Image Search functionality
package gis

import (
	"encoding/json"
	"errors"
	"gesture/core"
	"gesture/util"
	"math/rand"
	"net/url"
	"strings"
)

func Create(bot *core.Gobot) {
	bot.ListenFor("^gis (.*)", func(msg core.Message, matches []string) error {
		link, err := search(matches[1])
		if err == nil {
			msg.Ftfy(link)
		}
		return err
	})
}

// these structs really tie the room together, man
type gisResult struct {
	Url string
}
type gisResponse struct {
	ResponseData struct {
		Results []gisResult
	}
	Results []gisResult
}

// Search queries google for some images, and then randomly selects one
func search(search string) (result string, err error) {
	url := "http://ajax.googleapis.com/ajax/services/search/images?v=1.0&q=" + url.QueryEscape(search)
	body, err := util.GetUrl(url)
	if err != nil {
		return "", err
	}
	var gisResponse gisResponse
	json.Unmarshal(body, &gisResponse)
	if len(gisResponse.ResponseData.Results) > 0 {
		indexes := rand.Perm(len(gisResponse.ResponseData.Results))
		for _, index := range indexes {
			imageUrl := gisResponse.ResponseData.Results[index].Url
			contentType, err := util.ResponseHeaderContentType(imageUrl)
			if err == nil && strings.HasPrefix(contentType, "image/") {
				return imageUrl, nil
			}
		}
	}
	return "", errors.New("No image could be found for \"" + search + "\"")
}
