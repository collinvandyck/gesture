// does twitter related things
package twitter

import (
	"regexp"
	"strings"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

var (
	commandRegex = regexp.MustCompile(`^https?://(www.)?twitter.com/.*?/status/.*$`)
	httpClient = &http.Client{}
)

func IsStatusUrl(input string) bool {
	return commandRegex.MatchString(input)
}

// GetStatus queries a status url and outputs the rewritten text
func GetStatus(url string) (result string, err error) {
	pieces := strings.Split(url, "/")
	id := pieces[len(pieces)-1]
	statusUrl := "https://api.twitter.com/1/statuses/show/" + id + ".json"
	log.Printf("Getting status for url %s\n", statusUrl)
	resp, err := httpClient.Get(statusUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var jsonResponse map[string]interface{}
	json.Unmarshal(body, &jsonResponse)
	result = jsonResponse["text"].(string)
	return
}