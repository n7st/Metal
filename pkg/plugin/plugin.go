package plugin

import "github.com/n7st/metal/pkg/command"

// CommandFunction describes an individual plugin command function. These are
// matched against user-provided commands.
type CommandFunction func(*command.Command) *command.Response

// Plugin defines the interface for a bot plugin.
type Plugin interface {
	// Commands returns a list of the plugin's available commands against the
	// functions they must run.
	Commands() map[string]CommandFunction

	// Parse parses the input message
	Parse(*command.Command) *command.Response
}
