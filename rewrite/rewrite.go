/*
 Handles all of the rewriting tasks
*/
package rewrite

import (
	"github.com/collinvandyck/gesture/util"
	"log"
	"regexp"
	"strings"
)

var (
	linkPrefixes = []*regexp.Regexp{
		makeLinkRe("t.co"),
		makeLinkRe("cl.ly"),
		makeLinkRe("bit.ly"),
		makeLinkRe("j.mp"),
		makeLinkRe("tcrn.ch")}
	expanders       = []expander{expandUrl, expandEmbeddedImages}
	embeddedRePairs = []embeddedRePair{
		makeRePair(`(http://)?(www\.)?cl\.ly[^\s]+`, `a class="embed".*(http://cl\.ly[^"]+)`, 1),
		makeRePair(`(http://)?(www\.)?instagr.?am[^\s]+`, `img class="photo".*(http://[^"]+)`, 1),
		makeRePair(`(http://)?(x\.)?kingsh\.it[^\s]+`, `a class="embed".*(http://x\.kingsh\.it[^"]+)`, 1),
		makeRePair(`(https?://)?(www\.)?twitter\.com.*photo?[^\s]+`, `img src="(https?://[^"]+)".*Embedded image`, 1),
		makeRePair(`(https?://)?(www\.)?twitter\.com.*photo?[^\s]+`, `img.*media-slideshow-image.*src="(https?://[^"]+):.*".*`, 1),
	}
)

func makeLinkRe(part string) *regexp.Regexp {
	return regexp.MustCompile("^(http|https)?(://)?(www.)?" + part)
}

func makeRePair(link string, image string, imageSubmatch int) embeddedRePair {
	return embeddedRePair{
		regexp.MustCompile(link),
		regexp.MustCompile(image),
		imageSubmatch,
	}
}

type embeddedRePair struct {
	link          *regexp.Regexp // tests whether or not a token is a link
	image         *regexp.Regexp // what to search for in the fetched html
	imageSubmatch int            // what submatch to pull out of the image regexp
}

// GetRewrittenLinks takes an input line and rewrite any links that are shortened links into their full representation
// the return value is a slice of those rewritten links
func GetRewrittenLinks(input string) (result []string) {
	for _, link := range strings.Split(input, " ") {
		rewritten, err := expandAll(link)
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
		rewritten, err := expandAll(token)
		if err == nil && rewritten != "" {
			tokens[idx] = rewritten
		}
	}
	return strings.Join(tokens, " ")
}

// an expander is something that takes in a string and possibly expands it
type expander func(string) (string, error)

// thoroughly expand the input string by running it through the expander functions
func expandAll(input string) (string, error) {
	known := make(map[string]bool) // to track what we've seen already
	current := input
	known[current] = true
	for {
		rewritten := false
		for _, fn := range expanders {
			if result, err := fn(current); result != "" && err == nil {
				if known[result] {
					break
				}
				current = result
				known[current] = true
				rewritten = true
				break
			}
		}
		if !rewritten {
			break
		}
	}
	if current == input {
		return "", nil
	}
	return current, nil
}

func expandEmbeddedImages(url string) (result string, err error) {
	for _, rePair := range embeddedRePairs {
		if found := rePair.link.FindString(url); found != "" {
			body, err := util.GetUrl(found)
			if err != nil {
				return "", err
			}
			if matches := rePair.image.FindStringSubmatch(string(body)); matches != nil {
				return matches[rePair.imageSubmatch], nil
			}
		}
	}
	return "", nil
}

// expandUrl is an expander that expands shortened links
func expandUrl(url string) (result string, err error) {
	prefixFound := false
	for _, prefixRE := range linkPrefixes {
		if found := prefixRE.FindString(url); found != "" {
			prefixFound = true
			break
		}
	}
	if !prefixFound {
		return "", nil
	}
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	log.Println("Resolving url", url)
	return util.ResolveRedirects(url)
}
