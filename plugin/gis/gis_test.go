package gis

import (
	"strings"
	"testing"
)

func TestJpegSuffix(t *testing.T) {
	result := ensureSuffix("file.jpg", ".jpeg")
	if strings.HasSuffix(result, "jpeg") {
		t.Fatal("Should have jpg not jpeg")
	}
}

func TestSuffix(t *testing.T) {
	result := ensureSuffix("file", ".jpeg")
	if !strings.HasSuffix(result, ".jpeg") {
		t.Fatal("Should should have a .jpeg suffix")
	}
}
