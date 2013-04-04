// Google Image Search functionality
package gis

import (
	"errors"
	"gesture/core"
	"gesture/util"
	"math/rand"
	neturl "net/url"
	"strings"
	"time"
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
		return "", err
	}
	if len(gisResponse.ResponseData.Results) > 0 {

		// start a goroutine to determine image info for each response result
		imageUrlCh := make(chan string, len(gisResponse.ResponseData.Results))
		for _, resultUrl := range gisResponse.ResponseData.Results {
			go getImageInfo(resultUrl.Url, imageUrlCh)
		}

		// until a timeout is met, build a collection of urls
		urls := make([]string, 0, len(gisResponse.ResponseData.Results))
		timeout := time.After(500 * time.Millisecond)
		for {
			select {
			case url := <-imageUrlCh:
				urls = append(urls, url)
			case <-timeout:
				if len(urls) <= 0 {
					return "", errors.New("No image could be found for \"" + search + "\"")
				}
				indexes := rand.Perm(len(urls))
				for _, index := range indexes {
					url := urls[index]
					return url, nil
				}
			}
		}

	}
	return "", errors.New("No image could be found for \"" + search + "\"")
}

// getImageInfo looks at the header info for the url, and if it is an image, it sends an imageInfo on the channel
func getImageInfo(url string, ch chan string) {
	imageUrl, contentType, err := util.ResponseHeaderContentType(url)
	if err == nil && strings.HasPrefix(contentType, "image/") {
		select {
		case ch <- ensureSuffix(imageUrl, "."+contentType[len("image/"):]):
		default:
		}
	}
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
