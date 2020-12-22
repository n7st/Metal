package command

// Command defines a command to be passed to a plugin.
type Command struct {
	Message   string   // The original, unaltered message
	Argument  string   // The message less the first word as a string
	Arguments []string // The message less the first word (the command)
	Channel   string   // The channel the message originated in
	Command   string   // The first word of the message
	Username  string   // The name of the user who triggered the command
}

// Response defines a response to be sent to the chat system.
type Response struct {
	Message string // The message to output
	Channel string // The channel to send the message to
}
