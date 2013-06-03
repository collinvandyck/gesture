# What is this?

* gesture is an irc bot.
* gesture is an descendent of the allmighty <a href="http://github.com/dietrichf/jester">jester</a>.
* gesture runs on a plugin structure, still under construction, kinda like hubot.

# How do I plugin?

For better or worse, Plugins are just packages with a `func Create(bot
*core.Gobot)` method. There is no `interface` you should be defining. Within
that create method, a plugin is allowed to do anything it likes to the bot.
Check out the Gobot type for things you can do to your 'bot.

For example, this plugin inserts BEES into your IRC channel everytime someone
types the word "bees" (or "beeeees" or "beeeeeeeeees" and so on).
	
	# plugin/bees/bees.go

	package bees

	import (
		"github.com/collinvandyck/gesture/core"
	)

	func Create(bot *core.Gobot, config map[string]interface{}){
		bot.ListenFor("bee(e*)s", func(msg core.Message, matches []string) error {
			msg.Send("http://i.imgur.com/qrLEV.gif")
			return nil
		})
	}

# How do I configure?

Gesture uses the following struct for configuration:
 
    type Config struct {
         DisableFloodProtection bool
         BotName                string
         Hostname               string
         SSL                    bool
         Channels               []string
         Plugins                map[string]map[string]interface{}
    } 


Plugin configuration deserves a special mention. Plugin configuration is a map
of plugin names to plugin configuration maps. A plugin defines what it expects
out of a plugin configuration map (since it's a `map[string]interface{}`, feel
free to do whatever you please!). 

# Building your Robot


## Quickly

Gesture comes loaded with a `build_bot` tool that generates an IRC-bot script
for you based on the same JSON configuration file that your bot will load.

Install it with

    go install github.com/collinvandyck/gesture/build_bot

`build_bot` expects that your plugin names will be the full import path of the
plugin package, and will generate you a script that includes all of the plugins
listed in your configuration. Unfortunately, even plugins that don't need to be
configured should be listed.

An example configuration that you want to use `build_bot` on might look like:

	{
		"botname": "gesture",
		"hostname": "irc.freenode.net",
		"ssl": true,
		"channels": ["#lolgesture"],
		"plugins": {
      "github.com/collinvandyck/gesture/plugin/gis": {},
			"github.com/collinvandyck/gesture/plugin/graphite": {
				"prefix": "myprefix"
			},
			"github.com/collinvandyck/gesture/plugin/youtube": {
				"results": 3
			},
      "code.google.com/p/potato_chips_for_golang": {
        "flavor": "salt_and_vinegar"
      }
		}
	}

`build_bot` dumps it's results to standard output and leaves you to
save that to a file and modify it however you please. When you're ready to run
your bot, `go build` or `go run` it the usual way.

NOTE: `build_bot` does no validation of your config file *because* you have to
run it through the Go compiler yourself. If something's up with the output of
`build_bot`, check that you have the proper packages listed as the names of your
plugins.

## With Care

The `build_bot` script uses Go templates to do it's magic. Check out the script
for more detail on how to construct a template.

Please also feel free to ignore `build_bot` and write your own Gesture script!

