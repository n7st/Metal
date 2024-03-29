package main

import (
	"github.com/n7st/metal/pkg/command"
	"github.com/n7st/metal/pkg/plugin"
)

// ping contains the "ping" plugin.
type ping struct{}

func (p ping) Commands() map[string]plugin.CommandFunction {
	return map[string]plugin.CommandFunction{
		"ping": p.pingCommand,
	}
}

func (p ping) pingCommand(*command.Command) *command.Response {
	return &command.Response{Message: "Pong!"}
}

// NewPlugin initialises the "ping" plugin.
func NewPlugin() plugin.Plugin {
	return ping{}
}
