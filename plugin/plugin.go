package plugin

type Plugin interface {
	// allows the plugin process a received line
	Call(mc MessageContext) (bool, error)
}

// MessageContext is sent to each plugin's Call method
type MessageContext interface {
	Message() string
	Command() string
	CommandArgs() []string

	// this method replies to the original sender
	Reply(line string)

	// this method simply outputs some text to the channel
	Send(line string)
}

