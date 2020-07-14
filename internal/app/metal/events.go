// Package metal is the base for a pluggable IRC bot.
package metal

import (
	"strings"

	irc "github.com/thoj/go-ircevent"
)

// events returns all the bot's available IRC events.
func events(b *Bot) map[string]func(e *irc.Event) {
	return map[string]func(e *irc.Event){
		"001":     b.callback001,
		"PRIVMSG": b.callbackPrivmsg,
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

// callbackPrivmsg runs when the bot receives a message.
func (b *Bot) callbackPrivmsg(e *irc.Event) {
	args := strings.Split(e.Message(), " ")

	for _, module := range b.Modules {
		if module.Message {
			module.State.Global("Process")
			module.State.PushString(e.Message())
			module.State.Call(1, 1)

			v, ok := module.State.ToString(module.Next())

			if ok {
				b.Connection.Privmsg(e.Arguments[0], v)
			}
		}

		if len(args) >= 2 && module.Command {
			if command := module.Commands[args[0]]; command != "" {
				module.State.Global(command)
				module.State.PushString(args[1])
				module.State.Call(1, 1)

				v, ok := module.State.ToString(module.Next())

				if ok {
					b.Connection.Privmsg(e.Arguments[0], v)
				}
			}
		}
	}
}
