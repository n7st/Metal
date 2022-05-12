// Package metal is the base for a pluggable IRC bot.
package metal

import (
	"crypto/tls"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	irc "github.com/thoj/go-ircevent"

	"github.com/n7st/metal/internal/pkg/util"
)

// Bot contains the IRC bot.
type Bot struct {
	Connection *irc.Connection
	Config     *util.Config
	Logger     *logrus.Logger
	Plugins    []*util.Plugin

	Channels      map[string]bool
	ChannelsMutex *sync.RWMutex
}

// Init sets up an IRC bot connection to the network.
func Init(config *util.Config, logger *logrus.Logger) *Bot {
	connection := irc.IRC(config.IRC.Nickname, config.IRC.Ident)

	connection.Debug = config.IRC.Debug
	connection.VerboseCallbackHandler = config.IRC.Verbose
	connection.RealName = config.IRC.RealName

	if config.IRC.UseTLS {
		connection.UseTLS = true
		connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	err := connection.Connect(config.IRC.Hostname)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Fatal error connecting to IRC")
	}

	bot := &Bot{Connection: connection, Config: config, Logger: logger, ChannelsMutex: &sync.RWMutex{}}
	bot.Channels = make(map[string]bool)

	plugins, errors := util.LoadPlugins(config.EnabledPlugins())

	for _, p := range plugins {
		// Plugins may run on a ticker a goroutine; start them
		if timer, ok := p.Contrib.(interface{ Timer() }); ok {
			timer.Timer()
		}
	}

	if len(errors) > 0 {
		for _, err := range errors {
			logger.Println(err)
		}
	}

	bot.Plugins = plugins

	for name, fn := range events(bot) {
		bot.Connection.AddCallback(name, fn)
	}

	bot.healthCheck()

	return bot
}

// joinChannels joins either a provided list of channels, or the channels set in
// the bot's configuration.
func (b *Bot) joinChannels(params ...string) {
	var channels []string

	if len(params) > 0 {
		channels = params
	} else {
		channels = b.Config.IRC.Channels
	}

	for _, channel := range channels {
		b.Connection.Join(channel)
	}
}

func (b *Bot) IsOnChannel(channel string) bool {
	b.ChannelsMutex.RLock()

	defer b.ChannelsMutex.RUnlock()

	return b.Channels[channel]
}

// MessageChannels sends one message to many channels.
func (b *Bot) MessageChannels(channels []string, message string) {
	for _, channel := range channels {
		b.MessageChannel(channel, message)

		time.Sleep(1 * time.Second) // Antispam
	}
}

func (b *Bot) MessageChannel(channel string, message string) {
	if b.Config.IRC.Debug {
		b.Logger.Infof("Attempting to message channel %s with message %s", channel, message)
	}

	if b.IsOnChannel(channel) {
		b.Connection.Privmsg(channel, message)
	} else {
		b.Logger.Warnf("Tried to message channel %s with message %s, but I'm not in it", channel, message)
	}
}

// healthCheck checks for a connection to the IRC network and reconnects as
// required.
func (b *Bot) healthCheck() {
	retries := 0

	go func() {
		for {
			select {
			case <-b.Connection.Error:
				b.Logger.Warn("Healthcheck failed")

				if retries > b.Config.IRC.MaxReconnect {
					b.Logger.WithFields(logrus.Fields{
						"retries":     retries,
						"max_retries": b.Config.IRC.MaxReconnect,
					}).Fatal("Maximum reconnection attempts exceeded")
				}

				err := b.Connection.Reconnect()

				if err != nil {
					b.Logger.WithFields(logrus.Fields{
						"error": err.Error(),
					}).Warn("Health check error")
				} else {
					retries = 0
				}
			default:
				if b.Config.IRC.Verbose {
					b.Logger.Debug("Health check successful")
				}
			}

			time.Sleep(b.Config.IRC.ReconnectDelayMinutes * time.Minute)
		}
	}()
}
