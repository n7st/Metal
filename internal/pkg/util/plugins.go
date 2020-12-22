package util

import (
	loader "plugin"

	"github.com/n7st/metal/pkg/command"
	contrib "github.com/n7st/metal/pkg/plugin"
)

type Plugin struct {
	Contrib contrib.Plugin
}

func (p *Plugin) Run(c *command.Command) []*command.Response {
	var output []*command.Response

	for instruction, function := range p.Contrib.Commands() {
		if instruction == c.Command {
			response := function(c)

			if response != nil {
				output = append(output, response)
			}
		}
	}

	return output
}

func LoadPlugins() []*Plugin {
	script, err := loader.Open("plugins/binaries/greeter.so")

	if err != nil {
		panic(err)
	}

	initialiser, err := script.Lookup("NewPlugin")

	if err != nil {
		panic(err)
	}

	return []*Plugin{{Contrib: initialiser.(func() contrib.Plugin)()}}
}
