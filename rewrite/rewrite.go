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
	httpClient   = &http.Client{}
)

// GetRewrittenLinks takes an input line and rewrite any links that are shortened links into their full representation
// the return value is a slice of those rewritten links
func GetRewrittenLinks(input string) (result []string) {
	for _, link := range strings.Split(input, " ") {
		rewritten, err := expandUrl(link)
		if err == nil && rewritten != "" {
			result = append(result, rewritten)
		}
	}
	return
}

// Rewrite takes an input string, tokenizes it on whitespace, and then attempte to rewrite
// each token. The final result is joined back together at the end
func Rewrite(input string) string {
	tokens := strings.Split(input, " ")
	for idx, token := range tokens {
		rewritten, err := expandUrl(token)
		if err == nil && rewritten != "" {
			tokens[idx] = rewritten
		}
	}
	return strings.Join(tokens, " ")
}

// expandUrl is a rewriter that expands shortened links
func expandUrl(url string) (result string, err error) {
	for _, prefix := range linkPrefixes {
		if strings.HasPrefix(url, prefix) {
			break
		}
		return "", nil
	}

	if !strings.HasPrefix(url, "http") {
		log.Printf("Adding HTTP to url %s\n", url)
		url = "http://" + url
	}
	log.Printf("Expanding link %s\n", url)
	resp, err := httpClient.Head(url) // will follow redirects
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // not sure if i have to do this with a head response
	expanded := resp.Request.URL.String()
	if expanded != url {
		return expanded, nil
	}
	return "", nil
}
