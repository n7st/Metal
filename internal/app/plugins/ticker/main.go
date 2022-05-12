// ticker is an example of a plugin that runs a function on a schedule.
package main

import (
	"fmt"
	"time"

	"github.com/n7st/metal/pkg/plugin"
)

// ticker implements the plugin.Plugin interface, allowing it to be used as a
// bot plugin.
type ticker struct{}

func (t ticker) Timer() {
	go func() {
		for {
			fmt.Println("Tick")
			time.Sleep(time.Second * 10)
		}
	}()
}

// NewPlugin initialises the ticker plugin.
func NewPlugin() plugin.Plugin {
	return ticker{}
}
