/*
 The markov plugin listens passively to all messages and creates markov
 chains for each nick.
*/
package markov

import (
	"fmt"
	"github.com/collinvandyck/gesture/core"
	"github.com/collinvandyck/gesture/state"
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
	// todo: make prefix length configurable
	markov      = markovState{PrefixLength: 1, Chains: make(map[string]map[string][]string)}
	mutex       sync.Mutex
	pluginState = state.NewState("markov")
)

func Create(bot *core.Gobot, config map[string]interface{}) {
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
	bot.ListenFor("^ *markov +(.+)", func(msg core.Message, matches []string) core.Response {
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

// getChainMap gets the map for a particular user, or a new map with all of the data for all users
func getChainMap(user string) (map[string][]string, error) {
	if user != "" {
		userMap, ok := markov.Chains[user]
		if !ok {
			return nil, fmt.Errorf("No chain could be found for %s", user)
		}
		return userMap, nil
	}
	if len(markov.Chains) == 0 {
		return nil, fmt.Errorf("No chains could be found")
	}
	// combine all of the users' maps
	result := make(map[string][]string)
	for _, userChainMap := range markov.Chains {
		// userChainMap is a map[string][]string
		for prefix, userChain := range userChainMap {
			chain := result[prefix]
			if chain != nil {
				chain = make([]string,0)
			}
			for _, chainItem := range userChain {
				chain = append(chain, chainItem)
			}
			result[prefix] = chain
		}
	}
	return result, nil
}

func generateRandom() (string, error) {
	return generate("")
}

func generate(user string) (string, error) {
	chainMap, err := getChainMap(user)
	if err != nil {
		return "", err
	}
	p := newPrefix(markov.PrefixLength)
	var words []string
	for i := 0; i < maxWords; i++ {
		choices := chainMap[p.String()]
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
	return pluginState.Save(markov, false)
}

func contains(tokens []string, token string) bool {
	for _, word := range tokens {
		if word == token {
			return true
		}
	}
	return false
}
