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

	// this method simply outputs some text to the channel
	Send(line string)

	// sends a reply to the original sender
	Reply(line string)

	// sends a reply to the original sender with a ftfy prefix
	Ftfy(line string)

}
