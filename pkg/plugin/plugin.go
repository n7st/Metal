package plugin

import (
	"testing"

	"github.com/n7st/metal/pkg/command"
)

// CommandFunction describes an individual plugin command function. These are
// matched against user-provided commands.
type CommandFunction func(*command.Command) *command.Response

// Plugin defines the interface for a bot plugin. They may optionally implement
// three different functions:
// * Commands() map[string]plugin.CommandFunction - returns a list of commands
//   to the functions they trigger.
// * Parse(*command.Command) - for general purpose text processing (rather than
//   a specific command).
// * Timer() - for additional "timed" functions, e.g. a ticker in a goroutine.
type Plugin interface {
}

// CheckImplementsPluginInterface ensures a given plugin meets the Plugin
// interface.
func CheckImplementsInterface(t *testing.T, plugin interface{}) {
	var i interface{} = plugin

	if _, ok := i.(Plugin); !ok {
		t.Error("plugin does not implement the Plugin interface")
	}
}

// CheckRunCommand ensures a given plugin can run a named command. This is a
// helper for plugin tests.
func CheckRunCommand(t *testing.T, plugin Plugin, name string, input *command.Command) *command.Response {
	response := &command.Response{}

	if commander, ok := plugin.(interface {
		Commands() map[string]CommandFunction
	}); ok {
		function := commander.Commands()[name]

		if function == nil {
			t.Errorf("No such command available: %s", name)
		}

		response = function(input)
	}

	return response
}

// CheckRunParse ensures a given plugin can run the "Parse" function. This is a
// helper for plugin tests.
func CheckRunParse(t *testing.T, plugin Plugin, input *command.Command) *command.Response {
	response := &command.Response{}

	if parser, ok := plugin.(interface {
		Parse(c *command.Command) *command.Response
	}); ok {
		response = parser.Parse(input)
	}

	return response
}
