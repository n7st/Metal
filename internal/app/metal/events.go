// Package metal is the base for a pluggable IRC bot.
package metal

import (
	"strings"

	"github.com/n7st/metal/pkg/command"
	irc "github.com/thoj/go-ircevent"
)

// events returns all the bot's available IRC events.
func events(b *Bot) map[string]func(e *irc.Event) {
	return map[string]func(e *irc.Event){
		"001":     b.callback001,
		"privmsg": b.callbackPrivmsg,
	}
}

// callback001 runs when the bot connects to the network.
func (b *Bot) callback001(e *irc.Event) {
	if b.Config.IRC.Modes != "" {
		b.Connection.Mode(b.Connection.GetNick(), b.Config.IRC.Modes)
	}

	if b.Config.IRC.NickservPassword != "" {
		b.Connection.Privmsgf("nickserv", "identify %s %s",
			b.Config.IRC.NickservAccount,
			b.Config.IRC.NickservPassword,
		)
	}

	b.joinChannels()
}

func (b *Bot) callbackPrivmsg(e *irc.Event) {
	command := b.eventToCommand(e)

	for _, p := range b.Plugins {
		output := p.Run(command)

		for _, o := range output {
			b.Connection.Privmsg(e.Arguments[0], o.Message)
		}
	}
}

func (b *Bot) eventToCommand(e *irc.Event) *command.Command {
	msg := e.Message()
	instruction, arguments := b.messageToCommandArguments(msg)

	return &command.Command{
		Argument:  strings.Join(arguments, ""),
		Arguments: arguments,
		Command:   instruction,
		Message:   msg,
		Username:  e.Nick,
	}
}

func (b *Bot) messageToCommandArguments(message string) (string, []string) {
	var (
		arguments []string
		command   string
	)

	commandTrigger := b.Config.IRC.CommandTrigger

	if strings.HasPrefix(message, commandTrigger) {
		arguments = strings.Split(message, " ")
		command = strings.Replace(arguments[0], commandTrigger, "", 1)
		arguments = arguments[1:] // Clear the command from the arguments
	}

	return command, arguments
}
