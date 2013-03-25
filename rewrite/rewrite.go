/*
 Handles all of the rewriting tasks
*/
package rewrite

import (
	"log"
	"net/http"
	"strings"
	"gesture/twitter"
)

// a rewriter is something that can rewrite a string and if that happens, it will return that
// string, along with a true bool
type Rewriter func(string) (string, error)


var (
	linkPrefixes = []string{"t.co", "cl.ly", "www", "bit.ly", "j.mp", "tcrn.ch", "http"}
	httpClient   = &http.Client{}
)

// GetRewrittenLinks takes an input line and rewrite any links that are shortened links into their full representation
// the return value is a slice of those rewritten links
func GetRewrittenLinks(input string) (result []string) {
	for _, link := range strings.Split(input, " ") {
		rewritten, err := rewrite(link)
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
		rewritten, err := rewrite(token)
		if err == nil && rewritten != "" {
			tokens[idx] = rewritten
		}
	}
	return strings.Join(tokens, " ")
}

// the basic rewrite function. has a slice of Rewriters that it queries one by one. The first one
// that has a successful rewrite is the one that's used
func rewrite(token string) (result string, err error) {
	rewriters := []Rewriter{expandUrl, twitter.GetStatus}
	for _, rewriter := range rewriters {
		rewritten, err := rewriter(token)
		if err != nil {
			return "", err
		}
		if rewritten != "" {
			return rewritten, nil
		}				
	}
	return "", nil
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

