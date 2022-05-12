// Package util contains common functionality for "utilities" required by the
// bot.
//
// This file is for processing the application's configuration from a YAML file.
package util

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	configFilename  = "config.yaml" // These default values are used to create a
	vendorName      = "netsplit"    // filename for loading a file from the
	applicationName = "metal"       // platform's standard config location.

	defaultPort           = 6667
	defaultNickname       = "metalbot"
	defaultLogLevel       = "info"
	defaultCommandTrigger = "!"
	defaultMaxReconnect   = 3
)

// pluginConfig contains plugin-specific configuration.
type pluginConfig struct {
	Name    string
	Options map[string]interface{}
}

// ircConfig contains config items specific to the IRC bot itself.
type ircConfig struct {
	Channels         []string      `yaml:"channels"`
	CommandTrigger   string        `yaml:"command_trigger"`
	Debug            bool          `yaml:"debug"`
	Ident            string        `yaml:"ident"`
	MaxReconnect     int           `yaml:"max_reconnect"`
	ReconnectDelay   time.Duration `yaml:"reconnect_delay"`
	Modes            string        `yaml:"modes"`
	Nickname         string        `yaml:"nickname"`
	NickservAccount  string        `yaml:"nickserv_account"`
	NickservPassword string        `yaml:"nickserv_password"`
	Port             int           `yaml:"port"`
	RealName         string        `yaml:"real_name"`
	Server           string        `yaml:"server"`
	ServerPassword   string        `yaml:"server_password"`
	UseTLS           bool          `yaml:"use_tls"`
	Verbose          bool          `yaml:"verbose"`

	Hostname              string
	ReconnectDelayMinutes time.Duration
}

// Config contains the entire application's configuration.
type Config struct {
	IRC              *ircConfig               `yaml:"irc"`
	Plugins          map[string]*pluginConfig `yaml:"plugins"`
	UnparsedLogLevel string                   `yaml:"log_level"`

	LogLevel logrus.Level
}

// NewConfig sets up the application's configuration.
func NewConfig(params ...string) *Config {
	config := &Config{}
	config.Plugins = make(map[string]*pluginConfig)

	data, err := loadConfigData(params)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		panic(err)
	}

	config.applyDefaults()

	return config
}

// loadConfigData retrieves bytes from the config file. An optional filename can
// be provided to load configuration from a specific file rather than the
// default set for configdir.
func loadConfigData(params []string) ([]byte, error) {
	var (
		err      error
		data     []byte
		filename string
	)

	if len(params) > 0 {
		filename = params[0]

		data, err = ioutil.ReadFile(filename)
	} else {
		configDirs := configdir.New(vendorName, applicationName)
		folder := configDirs.QueryFolderContainsFile(configFilename)

		if folder != nil {
			data, err = folder.ReadFile(configFilename)
		}
	}

	if err != nil {
		panic(err)
	}

	return data, err
}

// applyDefaults sets default configuration values for items which are missing.
func (c *Config) applyDefaults() {
	if c.IRC == nil {
		c.IRC = &ircConfig{}
	}

	if c.IRC.Port == 0 {
		c.IRC.Port = defaultPort
	}

	if c.IRC.Nickname == "" {
		c.IRC.Nickname = defaultNickname
	}

	if c.IRC.Ident == "" {
		c.IRC.Ident = c.IRC.Nickname
	}

	if c.IRC.RealName == "" {
		c.IRC.RealName = c.IRC.Nickname
	}

	if c.IRC.NickservAccount == "" {
		c.IRC.NickservAccount = c.IRC.Nickname
	}

	if c.IRC.ReconnectDelay == 0 {
		c.IRC.ReconnectDelay = time.Duration(600 * time.Second)
	}

	if c.IRC.CommandTrigger == "" {
		c.IRC.CommandTrigger = defaultCommandTrigger
	}

	if c.IRC.MaxReconnect == 0 {
		c.IRC.MaxReconnect = defaultMaxReconnect
	}

	c.IRC.Hostname = fmt.Sprintf("%s:%d", c.IRC.Server, c.IRC.Port)

	c.setLogLevel()
}

// setLogLevel parses the configured logging level into one understood by
// logrus.
func (c *Config) setLogLevel() {
	level, err := logrus.ParseLevel(c.UnparsedLogLevel)

	if err != nil {
		level = logrus.InfoLevel
	}

	c.LogLevel = level
}

func (c *Config) EnabledPlugins() []string {
	var enabledPlugins []string

	// TODO: check for an "enabled" option, don't append if it exists and isn't
	// truthy
	for filepath := range c.Plugins {
		enabledPlugins = append(enabledPlugins, filepath)
	}

	return enabledPlugins
}
