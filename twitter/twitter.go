// does twitter related things
package twitter

import (
	"regexp"
	"strings"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"fmt"
)

var (
	commandRegex = regexp.MustCompile(`^https?://(www.)?twitter.com/.*?/status/.*$`)
	httpClient = &http.Client{}
)

func Describe(user string) (result string, err error) {
	url := "http://api.twitter.com/1/users/lookup.json?include_entities=true&screen_name=" + user
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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

// GetStatus queries a status url and outputs the rewritten text
func GetStatus(url string) (result string, err error) {
	if !commandRegex.MatchString(url) {
		return "", nil
	}
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