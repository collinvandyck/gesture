/*
 The markov plugin listens passively to all messages and creates markov
 chains for each nick.
*/
package markov

import (
	"fmt"
	"gesture/core"
	"gesture/state"
	"log"
	"math/rand"
	"strings"
	"sync"
)

type markovState struct {
	PrefixLength int
	Chains       map[string]map[string][]string // map[user]map[prefix][]chains
}

const (
	maxWords       = 100
	maxChainLength = 1000
)

var (
	markov      = markovState{PrefixLength: 2, Chains: make(map[string]map[string][]string)}
	mutex       sync.Mutex
	pluginState = state.NewState("markov")
)

func Create(bot *core.Gobot) {
	if err := pluginState.Load(&markov); err != nil {
		log.Printf("Could not load plugin state: %s", err)
	}

	bot.ListenFor("^ *markov *$", func(msg core.Message, matches []string) core.Response {
		mutex.Lock()
		defer mutex.Unlock()
		output, err := generateRandom()
		if err != nil {
			return bot.Error(err)
		}
		msg.Send(output)
		return bot.KeepGoing()
	})

	// generate a chain for the specified user
	bot.ListenFor("^ *markov *(.+)", func(msg core.Message, matches []string) core.Response {
		mutex.Lock()
		defer mutex.Unlock()
		output, err := generate(matches[1])
		if err != nil {
			return bot.Error(err)
		}
		msg.Send(output)
		return bot.KeepGoing()
	})

	// listen to everything
	bot.ListenFor("(.*)", func(msg core.Message, matches []string) core.Response {
		mutex.Lock()
		defer mutex.Unlock()
		user := msg.User
		text := matches[0]
		record(user, text)
		return bot.KeepGoing()
	})
}

func generateRandom() (string, error) {
	if len(markov.Chains) == 0 {
		return "", fmt.Errorf("No chains could be found")
	}
	users := make([]string, 0, len(markov.Chains))
	for k, _ := range markov.Chains {
		users = append(users, k)
	}
	user := users[rand.Intn(len(users))]
	// we have to unlock here b/c of deadlock caused by generate
	return generate(user)
}

func generate(user string) (string, error) {
	userMap, ok := markov.Chains[user]
	if !ok {
		return "", fmt.Errorf("No chain could be found for %s", user)
	}
	p := newPrefix(markov.PrefixLength)
	var words []string
	for i := 0; i < maxWords; i++ {
		choices := userMap[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.shift(next)
	}
	return strings.Join(words, " "), nil
}

// record breaks up the text into tokens and then creates chains for that user
func record(user, text string) error {
	p := newPrefix(markov.PrefixLength)
	tokens := strings.Split(text, " ")
	userMap, ok := markov.Chains[user]
	if !ok {
		markov.Chains[user] = make(map[string][]string)
		userMap = markov.Chains[user]
	}
	for _, token := range tokens {
		if strings.HasPrefix("http", token) {
			continue
		}
		str := p.String()
		if !contains(userMap[str], token) {
			userMap[str] = append(userMap[str], token)
			p.shift(token)
			// only allow maxChainLength items in a particular chain for a prefix
			if len(userMap[str]) > maxChainLength {
				userMap[str] = userMap[str][len(userMap[str])-maxChainLength:]
			}
		}
	}
	return pluginState.Save(markov)
}

func contains(tokens []string, token string) bool {
	for _, word := range tokens {
		if word == token {
			return true
		}
	}
	return false
}
