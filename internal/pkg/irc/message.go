// Package irc is a library for processing user input from IRC.
//
// This file converts inbound messages into usable commands and arguments.
package irc

import (
	"strings"

	irc "github.com/thoj/go-ircevent"
)

// Command is a command for an IRC bot. It can have arguments.
type Command struct {
	Argument  string
	Arguments []string
	Channel   string
	Command   string
	Message   string
}

// CommandConfig contains the application's available commands and the prefixes
// which trigger them.
//
// Prefixes tend to be individual characters (for example ! or .).
type CommandConfig struct {
	Commands []string
	Prefixes []string
}

// NewCommandConfig creates a CommandConfig.
func NewCommandConfig(commands []string, prefixes []string) *CommandConfig {
	return &CommandConfig{Commands: commands, Prefixes: prefixes}
}

// TryProcessCommand may return a Command{} if the user input is suitable.
// Suitability is determined based on the following:
//   - has a command prefix
//   - is a registered command
func (cc *CommandConfig) TryProcessCommand(e *irc.Event) *Command {
	if len(e.Arguments) < 1 {
		// I'm not sure we'll ever see an empty irc.Event, but it's better to be
		// safe.
		return nil
	}

	command := &Command{
		Message: e.Message(),
		Channel: e.Arguments[0],
	}

	if command.Channel == e.Connection.GetNick() {
		command.Channel = e.Nick // In the case of a direct message
	}

	parts := strings.Split(e.Message(), " ")

	if len(parts) > 0 {
		command.Command = parts[0]

		if len(parts) > 1 {
			command.Arguments = parts[1:]
			command.Argument = strings.Join(command.Arguments, " ")
		}
	}

	return command
}
