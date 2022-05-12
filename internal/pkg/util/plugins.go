package util

import (
	loader "plugin"

	"github.com/n7st/metal/pkg/command"
	contrib "github.com/n7st/metal/pkg/plugin"
)

// pluginFunctionName is the name of the function which returns the built
// plugin.
const pluginFunctionName = "NewPlugin"

// Plugin is a contributed plugin. It is wrapped with this struct so the Run
// function can be added for running user input against commands.
type Plugin struct {
	// Contrib is the external plugin.
	Contrib contrib.Plugin
}

// Run runs user input against the commands in the contributed plugin. Any
// command which returns input will be added to Run's output list.
func (p *Plugin) Run(c *command.Command) []*command.Response {
	var output []*command.Response

	if commander, ok := p.Contrib.(interface {
		Commands() map[string]contrib.CommandFunction
	}); ok {
		if c.Command != "" {
			for instruction, function := range commander.Commands() {
				if instruction == c.Command {
					if response := function(c); response != nil {
						output = append(output, response)
					}
				}
			}
		}
	}

	if parser, ok := p.Contrib.(interface {
		Parse(c *command.Command) *command.Response
	}); ok {
		if parsed := parser.Parse(c); parsed != nil {
			output = append(output, parsed)
		}
	}

	return output
}

// LoadPlugins loads contributed plugins by filename. Any plugin which fails to
// load or doesn't have the correct function will instead add an error to the
// errors list.
func LoadPlugins(filenames []string) ([]*Plugin, []error) {
	var (
		plugins []*Plugin
		errors  []error
	)

	for _, filename := range filenames {
		script, err := loader.Open(filename)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		initialiser, err := script.Lookup(pluginFunctionName)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		initialisedPlugin, err := initialisePlugin(initialiser)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		if initialisedPlugin != nil {
			plugins = append(plugins, &Plugin{Contrib: initialisedPlugin})
		}
	}

	return plugins, errors
}

func initialisePlugin(initialiser interface{}) (contrib.Plugin, error) {
	var (
		p   contrib.Plugin
		err error
	)

	defer func() {
		// A bad plugin should not cause the entire program to panic, which
		// would be the case if one was provided which didn't implement the
		// contrib.Plugin interface. Instead, catch the panic and return it as
		// a regular error.
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	p = initialiser.(func() contrib.Plugin)()

	return p, err
}
