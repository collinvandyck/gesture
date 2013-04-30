package markov

import (
	"strings"
)

type prefix struct {
	length int
	items  []string
}

func (p *prefix) String() string {
	return strings.Join(p.items, " ")
}

func newPrefix(length int) prefix {
	return prefix{length: length, items: make([]string, length)}
}

// shift moves the string into the rightmost slot, moving thens to the left
func (p *prefix) shift(token string) {
	copy(p.items, p.items[1:])
	p.items[len(p.items)-1] = token
}
