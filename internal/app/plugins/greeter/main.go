// greeter contains an example bot plugin with simple greeting commands.
package main

import (
	"strings"
	"sync"
	"time"

	"github.com/n7st/metal/pkg/command"
	"github.com/n7st/metal/pkg/plugin"
)

// greeter implements the plugin.Plugin interface, allowing it to be used as a
// bot plugin. All plugins must meet the definition to be loaded by the program.
type greeter struct {
	recentlyGreeted map[string]*time.Time
	rwm             *sync.RWMutex // For making the recentlyGreeted cache concurrency-safe
}

// Commands returns a list of the plugin's available commands to the functions
// they must run.
func (h greeter) Commands() map[string]plugin.CommandFunction {
	return map[string]plugin.CommandFunction{
		"hello_world": h.helloWorld,
	}
}

// Parse is a general-purpose function for processing messages in ways a command
// cannot. It can be used for text matching as seen here.
func (h greeter) Parse(c *command.Command) *command.Response {
	if strings.Contains(c.Message, "hi") {
		return &command.Response{
			Message: "Hi!",
		}
	}

	return nil
}

// NewPlugin initialises the greeter plugin.
func NewPlugin() plugin.Plugin {
	return greeter{
		recentlyGreeted: make(map[string]*time.Time),
		rwm:             &sync.RWMutex{},
	}
}

// helloWorld is an example "GreeterFunction". GreeterFunction is defined in the
// plugin module. It matches an individual command (in this case "hello_world")
// and returns a response.
func (h greeter) helloWorld(c *command.Command) *command.Response {
	if h.greetingPermitted(c.Username) {
		h.logGreeted(c.Username)

		who := "world"

		// c.Argument (string) and c.Arguments ([]string{}) can be used to
		// determine the output of your command
		if c.Argument != "" {
			who = c.Argument
		}

		return &command.Response{
			Message: "Hello, " + who,
		}
	}

	return &command.Response{
		Message: "You were greeted recently!",
	}
}

// greetingPermitted checks whether the user was greeted in the last minute.
func (h greeter) greetingPermitted(username string) bool {
	h.rwm.RLock()

	defer h.rwm.RUnlock()

	if lastGreeted := h.recentlyGreeted[username]; lastGreeted != nil {
		if diff := time.Now().Sub(*lastGreeted); diff >= time.Minute {
			return true
		}
	} else {
		return true
	}

	return false
}

// logGreeted stores the time at which a username was last greeted.
func (h greeter) logGreeted(username string) {
	h.rwm.RLock()

	defer h.rwm.RUnlock()

	greeted := time.Now()

	h.recentlyGreeted[username] = &greeted
}
