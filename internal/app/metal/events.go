// Package metal is the base for a pluggable IRC bot.
package metal

import irc "github.com/thoj/go-ircevent"

// events returns all the bot's available IRC events.
func events(b *Bot) map[string]func(e *irc.Event) {
	return map[string]func(e *irc.Event){
		"001": b.callback001,
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
