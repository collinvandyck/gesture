package core 

import (
	"errors"
)

// A Gobot plugin. Plugins are required to know their own name and to be able to
// create themselves given only a reference to the Gobot that they'll be
// plugged into.
type Plugin interface {
	Name() string
	Create(bot *Gobot) error
}

// A global plugin registry. Woof.
var registry = make(map[string]Plugin)

// Register a plugin. Only regsitered plugins will be visible to Gobots looking
// to extend themselves with awesome cool tricks.
//
// NOTE: Plugin names are unique. 
func Register(plugin Plugin) error {
	name := plugin.Name()

	if _, alreadyRegistered := registry[name]; alreadyRegistered {
		return errors.New("A plugin with name '" + name + "' already exists!")
	}
	registry[plugin.Name()] = plugin
	return nil
}

// Get a list all of the name of all available plugins
func ListPlugins() (names []string) {
	for name, _ := range registry {
		names = append(names, name)
	}
	return names
}

// Get all registered plugins
func GetAllPlugins() (plugins []Plugin) {
	for _, p := range registry {
		plugins = append(plugins, p)
	}
	return
}

// Return all plugins with the given names. Any plugin that can't be found
// generates an error.
func GetPlugins(names []string) (plugins []Plugin, es []error) {
	for _, name := range names {
		if p, exists := registry[name]; exists {
			plugins = append(plugins, p)
		} else {
			es = append(es, errors.New("No registered plugin named '"+name+"'"))
		}
	}
	return
}
