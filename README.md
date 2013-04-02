# What is this?

* gesture is an irc bot.
* is is an descendent of the allmighty <a href="http://github.com/dietrichf/jester">jester</a>.
* it runs on a plugin structure, still under construction, kinda like hubot.

# How do I plugin?


Plugins are required to know their own name and to be able to initialize themselves given a reference to the Gobot that they'll be plugged into.

```
type Plugin interface {
    Name() string
    Create(bot *Gobot) error
}
```

Gesture loads plugins from a global registry. To make your plugin available, call `core.Register` in your plugin file's `init()` function.

Here's an example of a simple plugin:
	
```go
// in plugin/bees.go
package plugin 

import (
	"gesture/core"
)

func init() {
	core.Register(Bees{})
}

type Bees struct{}

func (bees Bees) Name() string {
	return "bees"
}

func (bees Bees) Create(bot *core.Gobot) error {
	bot.ListenFor("bee(e*)s", func(msg core.Message, matches []string) error {
		msg.Send("http://i.imgur.com/qrLEV.gif")
		return nil
	})

	return nil
}
```



The Create(*core.Gobot) method allows the bot to register itself as a listener for particular regular expressions.  If something arrives on a channel that matches the regular expression, the plugin will be called with the Message along with any matching groups from the regular expression. The plugin is then able to reply back on that channel.

# How do I configure?

When starting gesture, the first argument to the program must be the location of a configuration file.  That configuration file should look something like this:

	{
		"botname": "gesturebot",
		"hostname": "irc.freenode.net",
		"ssl": true,
		"channels": ["#lolgesture"],
		"plugins": {
			"graphite": {
				"prefix": "myprefix"
			},
			"memegenerator": {
				"user": "foo",
				"password": "bar"
			},
			"youtube": {
				"results": 3
			}
		}
	}

Notice that each plugin, when created, will be passed a reference to the configuration so that it may use that data in its initialization process.

