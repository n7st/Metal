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
		"PRIVMSG": b.callbackPrivmsg,
		"JOIN":    b.callbackJoin,
		"KICK":    b.callbackKick,
		"PART":    b.callbackPart,
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
			if looksLikeChannel(e.Arguments[0]) {
				b.MessageChannel(e.Arguments[0], o.Message)
			} else {
				b.Connection.Privmsg(e.Arguments[0], o.Message)
			}
		}
	}
}

func (b *Bot) callbackJoin(e *irc.Event) {
	b.storeChannelStatus(e, true)
}

func (b *Bot) callbackPart(e *irc.Event) {
	b.storeChannelStatus(e, false)
}

func (b *Bot) callbackKick(e *irc.Event) {
	b.storeChannelStatus(e, false)
}

func (b *Bot) storeChannelStatus(e *irc.Event, onChannel bool) {
	b.ChannelsMutex.RLock()

	defer b.ChannelsMutex.RUnlock()

	channel := e.Arguments[0]

	if b.Config.IRC.Debug {
		b.Logger.Infof("Setting bot's channel presence for %s to %t", channel, onChannel)
	}

	if e.Nick == b.Connection.GetNick() {
		b.Channels[channel] = onChannel
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

func looksLikeChannel(input string) bool {
	return strings.HasPrefix(input, "#")
}
