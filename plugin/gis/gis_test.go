package gis

import (
	"strings"
	"testing"
)

func TestJpegSuffix(t *testing.T) {
	result, err := ensureSuffix("file.jpg", ".jpeg")
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasSuffix(result, "jpeg") {
		t.Fatal("Should have jpg not jpeg")
	}
}

func TestSuffix(t *testing.T) {
	result, err := ensureSuffix("file", ".jpeg")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(result, ".jpeg") {
		t.Fatal("Should should have a .jpeg suffix")
	}
}
