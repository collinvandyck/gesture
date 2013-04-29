// a helper package to let plugins and what not record and reload state
package state

import (
	"encoding/json"
	"os"
	"fmt"
	"bufio"
	"io/ioutil"
)

type State interface {
	// Save triggers a save of the state
	Save(data interface{}) error

	// Loads the data, possibly populating the container if no error
	Load(container interface{}) error
}

type state struct {
	name string
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

func (state *state) filename() string {
	return fmt.Sprintf("%s.json", state.name)
}

func (state *state) Save(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ioutil.WriteFile(state.filename(), bytes, 0666)
	return nil
}

// NewState builds a new State. The name will be used to save serialized data to {name}.state.
func NewState(name string) State {
	return &state{name}
}