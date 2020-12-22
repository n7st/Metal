package plugin

import (
	"testing"

	"github.com/n7st/metal/pkg/command"
)

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

// CheckImplementsPluginInterface ensures a given plugin meets the Plugin
// interface.
func CheckImplementsInterface(t *testing.T, plugin interface{}) {
	var i interface{} = plugin

	_, ok := i.(Plugin)

	if !ok {
		t.Error("plugin does not implement the Plugin interface")
	}
}

// CheckRunCommands ensures a given plugin can run a named command.
func CheckRunCommand(t *testing.T, plugin Plugin, name string, input *command.Command) *command.Response {
	function := plugin.Commands()[name]

	if function == nil {
		t.Errorf("No such command available: %s", name)
	}

	return function(input)
}
