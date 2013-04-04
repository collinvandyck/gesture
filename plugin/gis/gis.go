// Google Image Search functionality
package gis

import (
	"errors"
	"gesture/core"
	"gesture/util"
	"math/rand"
	neturl "net/url"
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
func search(search string) (string, error) {
	searchUrl := "http://ajax.googleapis.com/ajax/services/search/images?v=1.0&q=" + neturl.QueryEscape(search)
	var gisResponse gisResponse
	if err := util.UnmarshalUrl(searchUrl, &gisResponse); err != nil {
		return "", err;
	}
	if len(gisResponse.ResponseData.Results) > 0 {
		indexes := rand.Perm(len(gisResponse.ResponseData.Results))
		for _, index := range indexes {
			resultUrl := gisResponse.ResponseData.Results[index].Url
			imageUrl, contentType, err := util.ResponseHeaderContentType(resultUrl)
			if err == nil && strings.HasPrefix(contentType, "image/") {
				return ensureSuffix(imageUrl, "." + contentType[len("image/"):]), nil
			}
		}
	}
	return "", errors.New("No image could be found for \"" + search + "\"")
}

// ensureSuffix ensures a url ends with suffixes like .jpg, .png, etc
func ensureSuffix(url, suffix string) string {
	if strings.HasSuffix(strings.ToLower(url), strings.ToLower(suffix)) {
		return url
	}
	if strings.Contains(url, "?") {
		return url + "&lol" + suffix
	}
	return url + "?lol" + suffix
}

