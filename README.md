# What is this?

* gesture is an irc bot.
* is is an descendent of the allmighty <a href="http://github.com/dietrichf/jester">jester</a>.
* it runs on a plugin structure, still under construction, kinda like hubot.

# How do I plugin?

You will want to create a package in plugin/{packageName}.  For example,
	
	# plugin/bees/bees.go

	package bees

	import (
		"github.com/collinvandyck/gesture/core"
	)

	func Create(bot *core.Gobot){
		bot.ListenFor("bee(e*)s", func(msg core.Message, matches []string) error {
			msg.Send("http://i.imgur.com/qrLEV.gif")
			return nil
		})
	}


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

