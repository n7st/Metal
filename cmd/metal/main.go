// Package main contains the application's entrypoint.
package main

import (
	"github.com/sirupsen/logrus"

	"github.com/n7st/metal-core/internal/app/core"
	"github.com/n7st/metal-core/internal/pkg/util"
)

// main sets up an IRC bot and many RSS feed pollers.
func main() {
	config := util.NewConfig()
	logger := logrus.New()

	logger.SetLevel(config.LogLevel)
	logger.WithFields(logrus.Fields{
		"level": config.LogLevel.String(),
	}).Info("Set log level")

	bot := core.Init(config, logger)

	bot.Connection.Loop()
}
