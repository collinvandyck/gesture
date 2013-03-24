/*
 Handles all of the rewriting tasks
*/
package rewrite

import (
	"log"
	"net/http"
	"strings"
)

var (
	linkPrefixes = []string{"t.co", "cl.ly", "www", "bit.ly", "j.mp", "tcrn.ch", "http"}
	HttpClient   = &http.Client{}
)

// GetRewrittenLinks takes an input line and rewrite any links that are shortened links into their full representation
// the return value is a slice of those rewritten links
func GetRewrittenLinks(input string) (result []string) {
	for _, link := range findLinks(input) {
		expanded, _ := expandLink(link)
		if expanded != "" {
			result = append(result, expanded)
		}
	}
	return
}

// Rewrite takes an input string, tokenizes it on whitespace, and then attempte to rewrite
// each token. The final result is joined back together at the end
func Rewrite(input string) string {
	tokens := strings.Split(input, " ")
	for idx, token := range tokens {
		for _, prefix := range linkPrefixes {
			if strings.HasPrefix(token, prefix) {
				expanded, _ := expandLink(token)
				if expanded != "" {
					tokens[idx] = expanded
					break
				}
			}
		}
	}
	return strings.Join(tokens, " ")

}

// expandLink fully un-shortens a url
func expandLink(url string) (expanded string, err error) {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	log.Printf("Expanding link %s\n", url)
	resp, err := HttpClient.Head(url) // will follow redirects
	if err != nil {
		return expanded, err
	}
	defer resp.Body.Close() // not sure if i have to do this with a head response
	expanded = resp.Request.URL.String()
	if expanded != url {
		return
	}
	return "", nil
}

// findLinks returns a slice of strings that look like links. adds a protocol to the beginning of 
// the link if it doesn't already have one
func findLinks(message string) []string {
	result := make([]string, 0)
	for _, token := range strings.Split(message, " ") {
		// check to see if it looks like it might be a link
		for _, prefix := range linkPrefixes {
			if strings.HasPrefix(token, prefix) {
				result = append(result, "http://"+token)
				break
			}
		}
	}
	return result
}
