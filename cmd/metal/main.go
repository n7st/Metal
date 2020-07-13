// Package main contains the application's entrypoint.
package main

import (
	"github.com/sirupsen/logrus"

	"github.com/n7st/metal/internal/app/metal"
	"github.com/n7st/metal/internal/pkg/util"
)

// main sets up an IRC bot.
func main() {
	config := util.NewConfig()
	logger := logrus.New()

	logger.SetLevel(config.LogLevel)
	logger.WithFields(logrus.Fields{
		"level": config.LogLevel.String(),
	}).Info("Set log level")

	bot := metal.Init(config, logger)

	bot.Connection.Loop()
}
