// Configuration for your GoBot 
// "How to wave your hands"
package core 

type Config struct {
	BotName  string
	Hostname string
	SSL      bool
	Channels []string
	Plugins  map[string]map[string]interface{}
}
