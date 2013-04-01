// A Gesture interface to various YouTubery
package plugin 

import (
	"encoding/json"
	"errors"
	"fmt"
	"gesture/core"
	"gesture/util"
	"math/rand"
	"net/url"
	"regexp"
)

func init() {
	core.Register(YouTube{})
}

type YouTube struct {}

func (yt YouTube) Name() string {
	return "youtube"
}

func (yt YouTube) Create(bot *core.Gobot) (usingDefault error) {
	results, ok := bot.Config.Plugins["youtube"]["results"].(float64)
	if !ok {
		usingDefault = errors.New("Failed to load config for 'youtube' plugin. Using default result count of 1")
		results = 1
	}

	bot.ListenFor("^yt (.*)", func(msg core.Message, matches []string) error {
		link, err := yt.search(matches[1], int(results))
		if err == nil && link != "" {
			msg.Ftfy(link)
		}
		return err
	})

	return usingDefault
}


var urlCleaner = regexp.MustCompile(`&feature=youtube_gdata_player`)

// Search youtube for the given query string. Returns one of the first N youtube
// results for that search at random (everyone loves entropy!)
// Returns an empty string if there were no results for that query
func (yt YouTube) search(query string, results int) (link string, err error) {
	body, err := util.GetUrl(buildSearchUrl(query, results))
	if err != nil {
		return
	}

	var searchResponse youTubeResponse
	json.Unmarshal(body, &searchResponse)

	videos := searchResponse.Data.Items
	switch l := len(videos); {
	case l > 1:
		ordering := rand.Perm(len(videos))
		for _, i := range ordering {
			// Youtube adds a fragment to the end of players accessed via the API. Get
			// rid of that shit.
			link = urlCleaner.ReplaceAllLiteralString(videos[i].Player.Default, "")
		}
	case l == 1:
		link = urlCleaner.ReplaceAllLiteralString(videos[0].Player.Default, "")
	case l == 0:
		err = errors.New("No video found for search \"" + query + "\"")
	}

	return
}

// Generate a search URL for the given query. Returns the requested number of
// search results.
func buildSearchUrl(query string, results int) string {
	escapedQuery := url.QueryEscape(query)
	searchString := "https://gdata.youtube.com/feeds/api/videos?q=%v&max-results=%d&v=2&alt=jsonc"
	return fmt.Sprintf(searchString, escapedQuery, results)
}

// YouTube response types for deserializing JSON
type youTubePlayer struct {
	Default string
	Mobile  string
}

type youTubeItem struct {
	Title       string
	Description string
	Player      youTubePlayer
}

type youTubeData struct {
	Items []youTubeItem
}

type youTubeResponse struct {
	ApiVersion string
	Data       youTubeData
}
