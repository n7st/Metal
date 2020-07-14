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

	modules, loadErrors := util.LoadModules([]string{"modules/example.lua"})

	// Failed modules are non-fatal - they just won't ever be run
	for _, err := range loadErrors {
		logger.Warn(err.ToString())
	}

	logger.SetLevel(config.LogLevel)
	logger.WithFields(logrus.Fields{
		"level": config.LogLevel.String(),
	}).Info("Set log level")

	bot := metal.Init(config, logger, modules)

	bot.Connection.Loop()
}
