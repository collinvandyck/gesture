// A Gesture interface to various YouTubery
package youtube

import "fmt"
import "errors"
import "regexp"
import "strings"
import "net/url"
import "math/rand"
import "encoding/json"

import "gesture/util"
import "gesture/plugin"

// A YouTube plugin 
type YouTubePlugin struct {
	Results int
}

var urlCleaner = regexp.MustCompile(`&feature=youtube_gdata_player`)
var toob = YouTubePlugin{10} // Start with top 3 by relevance right now. Why not?

func NewPlugin() YouTubePlugin {
	return toob
}

// =============================================================================
// Plugin methods

func (plugin YouTubePlugin) Call(mc plugin.MessageContext) (success bool, err error) {
	success = false

	if mc.Command() == "yt" {
		if len(mc.CommandArgs()) > 0 {
			results, err := search(strings.Join(mc.CommandArgs(), " "), plugin.Results)
			if err != nil {
				return false, err
			}

			picked, err := pickRandomUrl(results)
			if err != nil {
				return false, err
			}
			picked = urlCleaner.ReplaceAllLiteralString(picked, " ")

			mc.Ftfy(picked)
		}
		success = true
	}
	return
}

// =============================================================================
// Functions

// Picks a random item from an array of youTubeItems and returns the url of the
// default player for that video.
func pickRandomUrl(videos []youTubeItem) (string, error) {
	if len(videos) > 0 {
		ordering := rand.Perm(len(videos))
		for _, i := range ordering {
			return videos[i].Player.Default, nil
		}
	}
	return "", errors.New("Can't sort an empty list!")
}

// Search youtube for the given query string. Returns the first N youTubeItems
// returned by the YouTube search API. 
func search(query string, results int) ([]youTubeItem, error) {
	body, err := util.GetUrl(buildSearchUrl(query, results))
	if err != nil {
		return nil, err
	}

	var searchResponse youTubeResponse
	json.Unmarshal(body, &searchResponse)

	return searchResponse.Data.Items, nil
}

// Generate a search URL for the given query. Returns the requested number of
// search results.
func buildSearchUrl(query string, results int) string {
	escapedQuery := url.QueryEscape(query)
	searchString := "https://gdata.youtube.com/feeds/api/videos?q=%v&max-results=%d&v=2&alt=jsonc"
	return fmt.Sprintf(searchString, escapedQuery, results)
}

// -----------------------------------------------------------------------------
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
	ApiVersion string      `json:"apiVersion"`
	Data       youTubeData `json:"data"`
}
