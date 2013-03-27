package rewrite

import (
	"regexp"
	"testing"
)

func TestCloudappLinks(t *testing.T) {
	url := "http://cl.ly/image/someimageid"
	re := regexp.MustCompile(`(http://)?(www\.)?cl\.ly`)
	found := re.Find([]byte(url))
	if found == nil {
		t.Error("Link did not match")
	}
}
