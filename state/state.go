// a helper package to let plugins and what not record and reload state
package state

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type State interface {
	// Save triggers a save of the state
	Save(data interface{}, force bool) error

	// Loads the data, possibly populating the container if no error
	Load(container interface{}) error
}

type state struct {
	name  string
	data  interface{} // the last data to be set
	mutex sync.Mutex
}

func (state *state) Load(container interface{}) error {
	file, err := os.Open(state.filename())
	if err != nil {
		return fmt.Errorf("Could not open file %s: %s", state.filename(), err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("Could not read file %s: %s", state.filename(), err)
	}
	err = json.Unmarshal(bytes, &container)
	if err != nil {
		return fmt.Errorf("Could not unmarshal json from file %s: %s", state.filename(), err)
	}
	return nil
}

// doSave does the actual saving. this is protected by a mutex
func (state *state) doSave(data interface{}) error {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	if data == nil {
		return nil
	}
	log.Printf("Saving %s", state.name)
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ioutil.WriteFile(state.filename(), bytes, 0666)
	// after we have successfully saved the state, set the data to nil so that it
	// is not re-saved on a subsequent doSave(...) without reason.
	state.data = nil
	return nil
}

func (state *state) filename() string {
	return fmt.Sprintf("%s.json", state.name)
}

// Save attempts to save the data. If force==true then this method is synchronous.
func (state *state) Save(data interface{}, force bool) error {
	if force {
		return state.doSave(data)
	}

	state.data = data
	return nil
}

// NewState builds a new State. The name will be used to save serialized data to {name}.state.
func NewState(name string) State {
	newState := &state{name: name}
	// spawn timed goroutine that will periodically save the state every minute
	go func() {
		for {
			timeout := time.After(1 * time.Minute)
			select {
			case <-timeout:
				if err := newState.doSave(newState.data); err != nil {
					log.Printf("Could not save %s on loop: %v", name, err)
				}
			}
		}
	}()
	return newState
}
