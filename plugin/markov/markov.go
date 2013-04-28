/*
 The markov plugin listens passively to all messages and creates markov
 chains for each nick.
*/
package markov

import (
	"fmt"
	"gesture/core"
	"strings"
	"sync"
	"math/rand"
)

type markovState struct {
	PrefixLength int
	Chains       map[string]map[string][]string // map[user]map[prefix][]chains	
}

const (
	maxWords = 100
)

var (
	state = markovState{PrefixLength: 2, Chains: make(map[string]map[string][]string)}
	mutex sync.Mutex
)

func Create(bot *core.Gobot) {
	// generate a chain for the specified user
	bot.ListenFor("markov *(.*)", func(msg core.Message, matches []string) error {
		output, err := generate(matches[1])
		if err != nil {
			return err
		}
		msg.Send(output)
		return nil
	})

	// listen to everything
	bot.ListenFor("(.*)", func(msg core.Message, matches []string) error {
		user := msg.User
		text := matches[0]
		record(user, text)
		return nil
	})
}

func generate(user string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	userMap, ok := state.Chains[user]
	if !ok {
		return "", fmt.Errorf("No chain could be found for %s", user)
	}
	p := newPrefix(state.PrefixLength)
	var words []string
	for i := 0; i < maxWords; i++ {
		choices := userMap[p.String()]		
		if len(choices) == 0 {
			break;
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.shift(next)
	}
	return strings.Join(words, " "), nil
}

// record breaks up the text into tokens and then creates chains for that user
func record(user, text string) {
	mutex.Lock()
	defer mutex.Unlock()

	p := newPrefix(state.PrefixLength)
	tokens := strings.Split(text, " ")
	userMap, ok := state.Chains[user]
	if !ok {
		state.Chains[user] = make(map[string][]string)
		userMap = state.Chains[user]
	}
	for _, token := range tokens {
		str := p.String()
		// todo: limit the length of the chain for this prefix
		userMap[str] = append(userMap[str], token)
		p.shift(token)
	}
}




