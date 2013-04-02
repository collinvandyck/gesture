package graphite

import (
	"testing"
)

func TestDereferenceNilMap(t *testing.T) {
	var m map[string]map[string]map[string]bool
	_, ok := m["foo"]["bar"]["baz"]
	if ok {
		t.Fatal("What?")
	}
}
